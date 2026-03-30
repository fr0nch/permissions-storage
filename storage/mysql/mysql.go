package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"permissions-storage/pkg/config"
	"permissions-storage/pkg/model"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fr0nch/logger"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db       *sql.DB
	settings *config.Settings
	log      *logger.Logger

	delayedQueryTimer *time.Timer
	cookies           []Cookie
	mu                sync.Mutex
}

type Cookie struct {
	userID model.UserID
	key    string
	value  any
}

func NewStorage(dsn url.URL, settings *config.Settings, logger *logger.Logger) (*Storage, error) {
	db, err := sql.Open("mysql", strings.TrimPrefix(dsn.String(), "mysql://"))
	if err != nil {
		logger.Errorf("Doesn't open database: %v\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(settings.Data.PoolSettings.MaximumPoolSize)
	db.SetMaxIdleConns(settings.Data.PoolSettings.MaximumIdle)
	db.SetConnMaxIdleTime(time.Duration(settings.Data.PoolSettings.MaximumIdleLifetime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(settings.Data.PoolSettings.MaximumLifetime) * time.Second)

	err = db.Ping()
	if err != nil {
		logger.Errorf("Doesn't open a connection: %v\n", err)
		return nil, err
	}

	logger.Debugf("Database connection established successfully")

	return &Storage{db: db, settings: settings, log: logger}, nil
}

func (p *Storage) Close() {
	p.db.Close()
}

func (p *Storage) CreateTables(ctx context.Context) error {
	//tx, err := p.db.BeginTx(ctx, nil)
	//if err != nil {
	//	return fmt.Errorf("failed to start the transaction.: %w", err)
	//}
	//
	//defer tx.Rollback()
	//
	//queries := []string{
	//	`CREATE TABLE IF NOT EXISTS groups (
	//        id int(11) NOT NULL AUTO_INCREMENT,
	//        name varchar(255) NOT NULL,
	//        priority int(11) NOT NULL DEFAULT 0,
	//        inheritance_id int(11) DEFAULT NULL,
	//        PRIMARY KEY (id),
	//        KEY fk_groups_groups (inheritance_id),
	//        CONSTRAINT fk_groups_groups FOREIGN KEY (inheritance_id) REFERENCES groups (id) ON DELETE SET NULL ON UPDATE NO ACTION
	//    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	//
	//	`CREATE TABLE IF NOT EXISTS group_options (
	//        id int(11) NOT NULL AUTO_INCREMENT,
	//        group_id int(11) NOT NULL,
	//        option_key varchar(255) NOT NULL,
	//        option_value longtext DEFAULT NULL,
	//        created_at timestamp NOT NULL DEFAULT current_timestamp(),
	//        updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	//        PRIMARY KEY (id) USING BTREE,
	//        UNIQUE KEY unique_group_option (group_id,option_key) USING BTREE,
	//        KEY group_id (group_id),
	//        KEY option_key (option_key),
	//        CONSTRAINT fk_group_options_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE NO ACTION
	//    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	//
	//	`CREATE TABLE IF NOT EXISTS group_permissions (
	//        id int(11) NOT NULL AUTO_INCREMENT,
	//        group_id int(11) NOT NULL,
	//        permission varchar(255) NOT NULL,
	//        PRIMARY KEY (id) USING BTREE,
	//        UNIQUE KEY unique_group_permission (group_id,permission),
	//        KEY group_id (group_id),
	//        KEY permission (permission),
	//        CONSTRAINT fk_group_permissions_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE NO ACTION
	//    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	//
	//	`CREATE TABLE IF NOT EXISTS users (
	//        ` + model.UserIDSQLColumnName + ` ` + model.UserIDSQLType + ` NOT NULL,
	//        name varchar(128) DEFAULT NULL,
	//        immunity int(11) NOT NULL DEFAULT 0,
	//        created_at timestamp NOT NULL DEFAULT current_timestamp(),
	//        lastvisit_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	//        PRIMARY KEY (steamid64)
	//    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	//
	//	`CREATE TABLE IF NOT EXISTS servers (
	//        id int(11) NOT NULL AUTO_INCREMENT,
	//        name varchar(255) NOT NULL,
	//        address varchar(32) DEFAULT NULL,
	//        default_group int(11) DEFAULT NULL,
	//        PRIMARY KEY (id),
	//        KEY fk_server_default_group (default_group),
	//        CONSTRAINT fk_server_default_group FOREIGN KEY (default_group) REFERENCES groups (id) ON DELETE SET NULL ON UPDATE NO ACTION
	//    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	//
	//	`CREATE TABLE IF NOT EXISTS server_groups (
	//        id int(11) NOT NULL AUTO_INCREMENT,
	//        server_id int(11) NOT NULL,
	//        group_id int(11) NOT NULL,
	//        PRIMARY KEY (id),
	//        UNIQUE KEY unique_server_groups (server_id,group_id) USING BTREE,
	//        KEY server_id (server_id),
	//        KEY group_id (group_id) USING BTREE,
	//        CONSTRAINT fk_server_groups_groups FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	//        CONSTRAINT fk_server_groups_servers FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE NO ACTION ON UPDATE NO ACTION
	//    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;`,
	//
	//	`CREATE TABLE IF NOT EXISTS server_user_groups (
	//        id int(11) NOT NULL AUTO_INCREMENT,
	//        ` + model.UserIDSQLColumnName + ` ` + model.UserIDSQLType + ` NOT NULL,
	//        server_id int(11) NOT NULL,
	//        group_id int(11) NOT NULL,
	//        created_at timestamp NOT NULL DEFAULT current_timestamp(),
	//        expires timestamp NULL DEFAULT '0000-00-00 00:00:00',
	//        updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	//        PRIMARY KEY (id),
	//        UNIQUE KEY unique_server_user_groups (server_id,` + model.UserIDSQLColumnName + `,group_id),
	//        KEY fk_server_user_groups_` + model.UserIDSQLColumnName + ` (` + model.UserIDSQLColumnName + `),
	//        KEY fk_server_user_groups_server_groups (server_id,group_id),
	//        CONSTRAINT fk_server_user_groups_server_groups FOREIGN KEY (server_id, group_id) REFERENCES server_groups (server_id, group_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	//        CONSTRAINT fk_server_user_groups_` + model.UserIDSQLColumnName + ` FOREIGN KEY (` + model.UserIDSQLColumnName + `) REFERENCES users (` + model.UserIDSQLColumnName + `) ON DELETE CASCADE ON UPDATE NO ACTION
	//    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	//
	//	`CREATE TABLE IF NOT EXISTS server_user_permissions (
	//        id int(11) NOT NULL AUTO_INCREMENT,
	//        ` + model.UserIDSQLColumnName + ` ` + model.UserIDSQLType + ` NOT NULL,
	//        server_id int(11) NOT NULL,
	//        permission varchar(255) NOT NULL,
	//        created_at timestamp NOT NULL DEFAULT current_timestamp(),
	//        expires timestamp NULL DEFAULT '0000-00-00 00:00:00',
	//        updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	//        PRIMARY KEY (id),
	//        UNIQUE KEY unique_server_user_permissions (` + model.UserIDSQLColumnName + `,server_id,permission) USING BTREE,
	//        KEY permission (permission),
	//        KEY server_id (server_id),
	//        KEY ` + model.UserIDSQLColumnName + ` (` + model.UserIDSQLColumnName + `),
	//        CONSTRAINT fk_user_permissions_server FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE ON UPDATE NO ACTION,
	//        CONSTRAINT fk_user_permissions_` + model.UserIDSQLColumnName + ` FOREIGN KEY (` + model.UserIDSQLColumnName + `) REFERENCES users (` + model.UserIDSQLColumnName + `) ON DELETE CASCADE ON UPDATE NO ACTION,
	//        CONSTRAINT chk_not_empty CHECK (trim(permission) <> '')
	//    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	//
	//	`CREATE TABLE IF NOT EXISTS user_cookies (
	//        id int(11) NOT NULL AUTO_INCREMENT,
	//        ` + model.UserIDSQLColumnName + ` ` + model.UserIDSQLType + ` NOT NULL,
	//        server_id int(11) NOT NULL,
	//        option_key varchar(255) NOT NULL,
	//        option_value longtext DEFAULT NULL,
	//        created_at timestamp NOT NULL DEFAULT current_timestamp(),
	//        updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	//        PRIMARY KEY (id),
	//        UNIQUE KEY unique_user_server_option (` + model.UserIDSQLColumnName + `,server_id,option_key),
	//        KEY fk_user_cookies_servers (server_id),
	//        KEY fk_user_cookies_` + model.UserIDSQLColumnName + ` (` + model.UserIDSQLColumnName + `),
	//        KEY option_key (option_key),
	//        CONSTRAINT fk_user_cookies_servers FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE ON UPDATE NO ACTION,
	//        CONSTRAINT fk_user_cookies_` + model.UserIDSQLColumnName + ` FOREIGN KEY (` + model.UserIDSQLColumnName + `) REFERENCES users (` + model.UserIDSQLColumnName + `) ON DELETE CASCADE ON UPDATE NO ACTION
	//    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	//}
	//
	//for _, query := range queries {
	//	_, err := tx.ExecContext(ctx, query)
	//	if err != nil {
	//		return fmt.Errorf("Query execution error:\n%s\n\nError: %w", query, err)
	//	}
	//}
	//
	//return tx.Commit()
	return nil
}

func (p *Storage) LoadGroups(ctx context.Context) (groups []*model.Group, defaultGroupID int, err error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, 0, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	groups, defaultID, err := p.loadGroups(ctx, tx, p.settings.ServerID)
	if err != nil {
		return nil, 0, err
	}

	if len(groups) == 0 {
		return nil, defaultID, tx.Commit()
	}

	ids := make([]int, len(groups))
	for i := range groups {
		ids[i] = groups[i].ID
		groups[i].Options = make(map[string]string)
	}

	var permsByGroup map[int][]string
	var optsByGroup map[int]map[string]string

	permsByGroup, err = p.loadGroupPermissions(ctx, tx, ids)
	if err != nil {
		return nil, 0, err
	}

	optsByGroup, err = p.loadGroupOptions(ctx, tx, ids)
	if err != nil {
		return nil, 0, err
	}

	for i := range groups {
		gid := groups[i].ID
		groups[i].Permissions = permsByGroup[gid]
		groups[i].Options = optsByGroup[gid]
	}

	return groups, defaultID, tx.Commit()
}

func (p *Storage) LoadUser(ctx context.Context, UserID model.UserID, username string) (*model.User, error) {
	user := &model.User{UserID: UserID}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `SELECT name, immunity FROM users WHERE steamid64 = ?`

	row := tx.QueryRowContext(ctx, query, UserID)
	if err := row.Scan(&user.Name, &user.Immunity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := p.addUser(ctx, tx, UserID, username); err != nil {
				return nil, fmt.Errorf("add user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("scan user: %w", err)
		}
	}

	user.Groups, err = p.loadUserGroups(ctx, tx, p.settings.ServerID, UserID)
	if err != nil {
		return nil, err
	}

	user.Permissions, err = p.loadUserPermissions(ctx, tx, p.settings.ServerID, UserID)
	if err != nil {
		return nil, err
	}

	var serverID int
	if !p.settings.GlobalCookie {
		serverID = p.settings.ServerID
	}

	user.Cookies, err = p.loadUserCookies(ctx, tx, serverID, UserID)
	if err != nil {
		return nil, err
	}

	return user, tx.Commit()
}

func (p *Storage) UpdateUser(ctx context.Context, user *model.User) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `UPDATE users SET name = ?, immunity = ?, lastvisit_at = ? WHERE steamid64 = ?;`

	_, err = tx.ExecContext(ctx, query, user.Name, user.Immunity, time.Now(), user.UserID)
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}

	return tx.Commit()
}

func (p *Storage) AddPermission(ctx context.Context, userID model.UserID, permission *model.UserPermission) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO server_user_permissions (steamid64, server_id, permission, expires)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
			expires = ?`

	_, err = tx.ExecContext(ctx, query, userID, p.settings.ServerID, permission.Permission, permission.Expires, permission.Expires)
	if err != nil {
		return fmt.Errorf("could not add permission: %w", err)
	}

	return tx.Commit()
}

func (p *Storage) RemovePermission(ctx context.Context, userID model.UserID, permission *model.UserPermission) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `DELETE FROM server_user_permissions WHERE steamid64 = ? AND server_id = ? AND permission = ?`

	_, err = tx.ExecContext(ctx, query, userID, p.settings.ServerID, permission.Permission)
	if err != nil {
		return fmt.Errorf("could not remove permission: %w", err)
	}

	return tx.Commit()
}

func (p *Storage) AddGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO server_user_groups (steamid64, server_id, group_id, expires)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
			expires = ?`

	_, err = tx.ExecContext(ctx, query, userID, p.settings.ServerID, group.GroupID, group.Expires, group.Expires)
	if err != nil {
		return fmt.Errorf("could not add group: %w", err)
	}

	return tx.Commit()
}

func (p *Storage) RemoveGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `DELETE FROM server_user_groups WHERE steamid64 = ? AND server_id = ? AND group_id = ?`

	_, err = tx.ExecContext(ctx, query, userID, p.settings.ServerID, group.GroupID)
	if err != nil {
		return fmt.Errorf("could not remove group: %w", err)
	}

	return tx.Commit()
}

func (p *Storage) SetCookie(userID model.UserID, key string, value any) {
	p.addQuery(Cookie{
		userID: userID,
		key:    key,
		value:  value,
	})
}

func (p *Storage) addQuery(cookie Cookie) {
	p.mu.Lock()

	index := slices.IndexFunc(p.cookies, func(c Cookie) bool {
		return c.key == cookie.key
	})

	if index != -1 {
		p.cookies[index] = cookie
	} else {
		p.cookies = append(p.cookies, cookie)
	}

	if p.delayedQueryTimer != nil {
		p.delayedQueryTimer.Reset(5 * time.Second)
		p.mu.Unlock()
		return
	}

	p.delayedQueryTimer = time.AfterFunc(5*time.Second, func() {
		tmpCookies := make([]Cookie, len(p.cookies))

		p.mu.Lock()
		copy(tmpCookies, p.cookies)
		clear(p.cookies)
		p.mu.Unlock()

		for _, _cookie := range tmpCookies {
			err := p.setCookie(context.Background(), _cookie.userID, _cookie.key, _cookie.value)
			if err != nil {
				p.log.Errorf("Error setting cookie: %v\n", err)
			}
		}
	})

	p.mu.Unlock()
}

func (p *Storage) setCookie(ctx context.Context, userID model.UserID, key string, value any) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO user_cookies (steamid64, server_id, option_key, option_value)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
			option_value = ?`

	var serverID int
	if !p.settings.GlobalCookie {
		serverID = p.settings.ServerID
	}

	_, err = tx.ExecContext(ctx, query, userID, serverID, key, value, value)
	if err != nil {
		return fmt.Errorf("could not setup cookie: %w", err)
	}

	return tx.Commit()
}

func (p *Storage) loadGroups(ctx context.Context, tx *sql.Tx, serverID int) ([]*model.Group, int, error) {
	const query = `
		SELECT  g.id,
				g.name,
				g.priority,
				g.inheritance_id,
				(COALESCE(s.default_group, s0.default_group, 0) = g.id) AS is_default
		FROM server_groups sg
		JOIN groups g ON g.id = sg.group_id
		LEFT JOIN servers s ON s.id = ?
		LEFT JOIN servers s0 ON s0.id = 0
		WHERE sg.server_id IN (?, 0)
		ORDER BY g.priority DESC;
	`

	rows, err := tx.QueryContext(ctx, query, serverID, serverID)
	if err != nil {
		return nil, 0, fmt.Errorf("query groups: %w", err)
	}

	defer rows.Close()

	res := make([]*model.Group, 0)
	var defaultID int

	for rows.Next() {
		g := &model.Group{}
		if err := rows.Scan(&g.ID, &g.Name, &g.Priority, &g.InheritanceID, &g.Default); err != nil {
			return nil, 0, fmt.Errorf("scan group: %w", err)
		}
		if g.Default {
			defaultID = g.ID
		}
		res = append(res, g)
	}

	return res, defaultID, rows.Err()
}

func (p *Storage) loadGroupPermissions(ctx context.Context, tx *sql.Tx, groupIDs []int) (map[int][]string, error) {
	in, args := makeInClause(groupIDs)

	query := fmt.Sprintf(`
		SELECT group_id, permission
		FROM group_permissions
		WHERE group_id IN (%s)
		ORDER BY group_id;
	`, in)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query manager: %w", err)
	}

	defer rows.Close()

	out := make(map[int][]string, len(groupIDs))
	for rows.Next() {
		var gid int
		var perm string

		if err := rows.Scan(&gid, &perm); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}

		out[gid] = append(out[gid], perm)
	}

	return out, rows.Err()
}

func (p *Storage) loadGroupOptions(ctx context.Context, tx *sql.Tx, groupIDs []int) (map[int]map[string]string, error) {
	in, args := makeInClause(groupIDs)

	query := fmt.Sprintf(`
		SELECT group_id, option_key, option_value
		FROM group_options
		WHERE group_id IN (%s)
		ORDER BY group_id;
	`, in)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query options: %w", err)
	}
	defer rows.Close()

	out := make(map[int]map[string]string, len(groupIDs))
	for rows.Next() {
		var gid int
		var k, v string

		if err := rows.Scan(&gid, &k, &v); err != nil {
			return nil, fmt.Errorf("scan option: %w", err)
		}

		if out[gid] == nil {
			out[gid] = make(map[string]string)
		}

		out[gid][k] = v
	}

	return out, rows.Err()
}

func (p *Storage) addUser(ctx context.Context, tx *sql.Tx, UserID model.UserID, username string) error {
	query := `INSERT INTO users (steamid64, name, immunity) VALUES (?, ?, ?)`

	_, err := tx.ExecContext(ctx, query, UserID, username, 0)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (p *Storage) loadUserGroups(ctx context.Context, tx *sql.Tx, serverID int, UserID model.UserID) ([]model.UserGroup, error) {
	query := `
		SELECT
			group_id,
			g.name,
			expires
		FROM server_user_groups
		INNER JOIN groups g ON group_id = g.id
		WHERE steamid64 = ? AND server_id = ? AND (expires IS NULL OR expires = 0 OR expires > CURRENT_TIMESTAMP())
	`

	rows, err := tx.QueryContext(ctx, query, UserID, serverID)
	if err != nil {
		return nil, fmt.Errorf("query options: %w", err)
	}

	defer rows.Close()

	var groups = make([]model.UserGroup, 0)

	for rows.Next() {
		var group model.UserGroup
		var expires sql.NullTime

		if err := rows.Scan(&group.GroupID, &group.GroupName, &expires); err != nil {
			p.log.Errorf("Failed to scan server user groups: %v\n", err)
			continue
		}

		if expires.Valid {
			group.Expires = expires.Time
		} else {
			group.Expires = time.Time{}
		}

		p.log.Debugf("Loading group '%s'[id: %d] for user: %d\n", group.GroupName, group.GroupID, UserID)
		groups = append(groups, group)
	}

	return groups, rows.Err()
}

func (p *Storage) loadUserPermissions(ctx context.Context, tx *sql.Tx, serverID int, UserID model.UserID) ([]model.UserPermission, error) {
	query := `
		SELECT
			permission,
			expires
		FROM server_user_permissions 
		WHERE steamid64 = ? AND server_id = ? AND (expires IS NULL OR expires = 0) OR (expires > CURRENT_TIMESTAMP())
	`

	rows, err := tx.QueryContext(ctx, query, UserID, serverID)
	if err != nil {
		return nil, fmt.Errorf("query options: %w", err)
	}

	defer rows.Close()

	var permissions = make([]model.UserPermission, 0)

	for rows.Next() {
		var permission model.UserPermission
		var expires sql.NullTime

		if err := rows.Scan(&permission.Permission, &expires); err != nil {
			p.log.Errorf("Failed to scan server user manager: %v\n", err)
			continue
		}

		if expires.Valid {
			permission.Expires = expires.Time
		} else {
			permission.Expires = time.Time{}
		}

		p.log.Debugf("Loading permission '%s' for user: %d\n", permission.Permission, UserID)
		permissions = append(permissions, permission)
	}

	p.log.Debugf("Loaded %d permissions\n", len(permissions))

	return permissions, rows.Err()
}

func (p *Storage) loadUserCookies(ctx context.Context, tx *sql.Tx, serverID int, UserID model.UserID) (map[string]string, error) {
	query := `
		SELECT option_key, option_value 
		FROM user_cookies 
		WHERE steamid64 = ? AND server_id = ?
	`

	rows, err := tx.QueryContext(ctx, query, UserID, serverID)
	if err != nil {
		return nil, fmt.Errorf("query options: %w", err)
	}

	defer rows.Close()

	var cookies = make(map[string]string)

	for rows.Next() {
		var cookieKey, cookieValue string
		if err := rows.Scan(&cookieKey, &cookieValue); err != nil {
			p.log.Errorf("Failed to scan server user cookies: %v\n", err)
			continue
		}

		p.log.Debugf("Loading cookie '%s' for user: %d\n", cookieKey, UserID)
		cookies[cookieKey] = cookieValue
	}

	p.log.Debugf("Loaded %d cookies\n", len(cookies))

	return cookies, rows.Err()
}

func makeInClause(ids []int) (string, []any) {
	parts := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))

	for _, id := range ids {
		parts = append(parts, "?")
		args = append(args, id)
	}

	return strings.Join(parts, ", "), args
}
