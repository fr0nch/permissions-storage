package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"permissions-storage/pkg/config"
	"permissions-storage/pkg/model"
	"slices"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/fr0nch/logger"
)

type Storage struct {
	db *sql.DB

	ready  chan struct{}
	failed chan struct{}

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
	s := &Storage{
		ready:    make(chan struct{}),
		failed:   make(chan struct{}),
		settings: settings,
		log:      logger,
	}

	dir := filepath.Dir(dsn.Path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	db, err := sql.Open("sqlite3", dsn.String())
	if err != nil {
		logger.Errorf("Doesn't open database: %v\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(settings.Data.PoolSettings.MaximumPoolSize)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		close(s.failed)

		return nil, fmt.Errorf("database is unreachable: %w", err)
	}

	s.db = db

	logger.Debugf("Database connection established successfully")
	close(s.ready)

	return s, nil
}

func (s *Storage) WaitReady(ctx context.Context) error {
	select {
	case <-s.ready:
		return nil
	case <-s.failed:
		return fmt.Errorf("database connection failed")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) CreateTables(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start the transaction.: %w", err)
	}

	defer tx.Rollback()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			priority INTEGER NOT NULL DEFAULT 0,
			inheritance_id INTEGER DEFAULT NULL,
			FOREIGN KEY (inheritance_id) REFERENCES groups (id) ON DELETE SET NULL
		);`,

		`CREATE INDEX IF NOT EXISTS idx_groups_inheritance_id ON groups(inheritance_id);`,

		`CREATE TABLE IF NOT EXISTS group_options (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			option_key TEXT NOT NULL,
			option_value TEXT DEFAULT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(group_id, option_key),
			FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
		);`,

		`CREATE INDEX IF NOT EXISTS idx_group_options_group_id ON group_options(group_id);`,
		`CREATE INDEX IF NOT EXISTS idx_group_options_key ON group_options(option_key);`,

		`CREATE TABLE IF NOT EXISTS group_permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			permission TEXT NOT NULL,
			UNIQUE(group_id, permission),
			FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
		);`,

		`CREATE INDEX IF NOT EXISTS idx_group_permissions_group_id ON group_permissions(group_id);`,
		`CREATE INDEX IF NOT EXISTS idx_group_permissions_permission ON group_permissions(permission);`,

		`CREATE TABLE IF NOT EXISTS users (
			steamid64 INTEGER PRIMARY KEY,
			name TEXT DEFAULT NULL,
			immunity INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			lastvisit_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS servers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			address TEXT DEFAULT NULL,
			default_group INTEGER DEFAULT NULL,
			FOREIGN KEY (default_group) REFERENCES groups (id) ON DELETE SET NULL
		);`,

		`CREATE INDEX IF NOT EXISTS idx_servers_default_group ON servers(default_group);`,

		`CREATE TABLE IF NOT EXISTS server_groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			server_id INTEGER NOT NULL,
			group_id INTEGER NOT NULL,
			UNIQUE(server_id, group_id),
			FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE NO ACTION,
			FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE NO ACTION
		);`,

		`CREATE INDEX IF NOT EXISTS idx_server_groups_server_id ON server_groups(server_id);`,
		`CREATE INDEX IF NOT EXISTS idx_server_groups_group_id ON server_groups(group_id);`,

		`INSERT OR IGNORE INTO groups (id, name, priority, inheritance_id) VALUES (1, 'Default', 0, NULL);`,
		`INSERT OR IGNORE INTO servers (id, name, address, default_group) VALUES (0, 'All Servers', NULL, 1);`,
		`INSERT OR IGNORE INTO server_groups (id, server_id, group_id) VALUES (1, 0, 1);`,

		`CREATE TABLE IF NOT EXISTS server_user_groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			steamid64 INTEGER NOT NULL,
			server_id INTEGER NOT NULL,
			group_id INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires TIMESTAMP DEFAULT NULL,
			updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(server_id, steamid64, group_id),
			FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE,
			FOREIGN KEY (server_id, group_id) REFERENCES server_groups (server_id, group_id) ON DELETE CASCADE
		);`,

		`CREATE INDEX IF NOT EXISTS idx_server_user_groups_steamid64 ON server_user_groups(steamid64);`,
		`CREATE INDEX IF NOT EXISTS idx_server_user_groups_server_group ON server_user_groups(server_id, group_id);`,

		`CREATE TABLE IF NOT EXISTS server_user_permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			steamid64 INTEGER NOT NULL,
			server_id INTEGER NOT NULL,
			permission TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires TIMESTAMP DEFAULT NULL,
			updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(steamid64, server_id, permission),
			FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE,
			FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE,
			CHECK (TRIM(permission) != '')
		);`,

		`CREATE INDEX IF NOT EXISTS idx_server_user_permissions_permission ON server_user_permissions(permission);`,
		`CREATE INDEX IF NOT EXISTS idx_server_user_permissions_server_id ON server_user_permissions(server_id);`,
		`CREATE INDEX IF NOT EXISTS idx_server_user_permissions_steamid64 ON server_user_permissions(steamid64);`,

		`CREATE TABLE IF NOT EXISTS user_cookies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			steamid64 INTEGER NOT NULL,
			server_id INTEGER NOT NULL,
			option_key TEXT NOT NULL,
			option_value TEXT DEFAULT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(steamid64, server_id, option_key),
			FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE,
			FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE
		);`,

		`CREATE INDEX IF NOT EXISTS idx_user_cookies_server_id ON user_cookies(server_id);`,
		`CREATE INDEX IF NOT EXISTS idx_user_cookies_steamid64 ON user_cookies(steamid64);`,
		`CREATE INDEX IF NOT EXISTS idx_user_cookies_option_key ON user_cookies(option_key);`,
	}

	for _, query := range queries {
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("Query execution error:\n%s\n\nError: %w", query, err)
		}
	}

	return tx.Commit()
}

func (s *Storage) LoadGroups(ctx context.Context) (groups []*model.Group, defaultGroupID int, err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, 0, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	groups, defaultID, err := s.loadGroups(ctx, tx, s.settings.ServerID)
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

	permsByGroup, err = s.loadGroupPermissions(ctx, tx, ids)
	if err != nil {
		return nil, 0, err
	}

	optsByGroup, err = s.loadGroupOptions(ctx, tx, ids)
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

func (s *Storage) LoadUser(ctx context.Context, UserID model.UserID, username string) (*model.User, error) {
	user := &model.User{UserID: UserID}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `SELECT name, immunity FROM users WHERE steamid64 = ?`

	row := tx.QueryRowContext(ctx, query, UserID)
	if err := row.Scan(&user.Name, &user.Immunity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := s.addUser(ctx, tx, UserID, username); err != nil {
				return nil, fmt.Errorf("add user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("scan user: %w", err)
		}
	}

	user.Groups, err = s.loadUserGroups(ctx, tx, s.settings.ServerID, UserID)
	if err != nil {
		return nil, err
	}

	user.Permissions, err = s.loadUserPermissions(ctx, tx, s.settings.ServerID, UserID)
	if err != nil {
		return nil, err
	}

	var serverID int
	if !s.settings.GlobalCookie {
		serverID = s.settings.ServerID
	}

	user.Cookies, err = s.loadUserCookies(ctx, tx, serverID, UserID)
	if err != nil {
		return nil, err
	}

	return user, tx.Commit()
}

func (s *Storage) UpdateUser(ctx context.Context, user *model.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
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

func (s *Storage) AddPermission(ctx context.Context, userID model.UserID, permission *model.UserPermission) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO server_user_permissions (steamid64, server_id, permission, expires)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(steamid64, server_id, permission)
			DO UPDATE SET expires = ?, updated = CURRENT_TIMESTAMP`

	_, err = tx.ExecContext(ctx, query, userID, s.settings.ServerID, permission.Permission, expiry(permission.Expires), expiry(permission.Expires))
	if err != nil {
		return fmt.Errorf("could not add permission: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) RemovePermission(ctx context.Context, userID model.UserID, permission *model.UserPermission) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `DELETE FROM server_user_permissions WHERE steamid64 = ? AND server_id = ? AND permission = ?`

	_, err = tx.ExecContext(ctx, query, userID, s.settings.ServerID, permission.Permission)
	if err != nil {
		return fmt.Errorf("could not remove permission: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) AddGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO server_user_groups (steamid64, server_id, group_id, expires)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(steamid64, server_id, group_id)
			DO UPDATE SET expires = ?, updated = CURRENT_TIMESTAMP`

	_, err = tx.ExecContext(ctx, query, userID, s.settings.ServerID, group.GroupID, expiry(group.Expires), expiry(group.Expires))
	if err != nil {
		return fmt.Errorf("could not add group: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) RemoveGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `DELETE FROM server_user_groups WHERE steamid64 = ? AND server_id = ? AND group_id = ?`

	_, err = tx.ExecContext(ctx, query, userID, s.settings.ServerID, group.GroupID)
	if err != nil {
		return fmt.Errorf("could not remove group: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) SetCookie(userID model.UserID, key string, value any) {
	s.addQuery(Cookie{
		userID: userID,
		key:    key,
		value:  value,
	})
}

func (s *Storage) addQuery(cookie Cookie) {
	s.mu.Lock()

	index := slices.IndexFunc(s.cookies, func(c Cookie) bool {
		return c.key == cookie.key
	})

	if index != -1 {
		s.cookies[index] = cookie
	} else {
		s.cookies = append(s.cookies, cookie)
	}

	if s.delayedQueryTimer != nil {
		s.delayedQueryTimer.Reset(5 * time.Second)
		s.mu.Unlock()
		return
	}

	s.delayedQueryTimer = time.AfterFunc(5*time.Second, func() {
		tmpCookies := make([]Cookie, len(s.cookies))

		s.mu.Lock()
		copy(tmpCookies, s.cookies)
		clear(s.cookies)
		s.mu.Unlock()

		for _, _cookie := range tmpCookies {
			err := s.setCookie(context.Background(), _cookie.userID, _cookie.key, _cookie.value)
			if err != nil {
				s.log.Errorf("Error setting cookie: %v\n", err)
			}
		}
	})

	s.mu.Unlock()
}

func (s *Storage) setCookie(ctx context.Context, userID model.UserID, key string, value any) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO user_cookies (steamid64, server_id, option_key, option_value)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(steamid64, server_id, option_key)
			DO UPDATE SET option_value = ?, updated = CURRENT_TIMESTAMP`

	var serverID int
	if !s.settings.GlobalCookie {
		serverID = s.settings.ServerID
	}

	_, err = tx.ExecContext(ctx, query, userID, serverID, key, value, value)
	if err != nil {
		return fmt.Errorf("could not setup cookie: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) loadGroups(ctx context.Context, tx *sql.Tx, serverID int) ([]*model.Group, int, error) {
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

func (s *Storage) loadGroupPermissions(ctx context.Context, tx *sql.Tx, groupIDs []int) (map[int][]string, error) {
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

func (s *Storage) loadGroupOptions(ctx context.Context, tx *sql.Tx, groupIDs []int) (map[int]map[string]string, error) {
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

func (s *Storage) addUser(ctx context.Context, tx *sql.Tx, UserID model.UserID, username string) error {
	query := `INSERT INTO users (steamid64, name, immunity) VALUES (?, ?, ?)`

	_, err := tx.ExecContext(ctx, query, UserID, username, 0)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (s *Storage) loadUserGroups(ctx context.Context, tx *sql.Tx, serverID int, UserID model.UserID) ([]model.UserGroup, error) {
	query := `
		SELECT
			group_id,
			g.name,
			expires
		FROM server_user_groups
		INNER JOIN groups g ON group_id = g.id
		WHERE steamid64 = ? AND server_id = ? AND (expires IS NULL OR expires = 0 OR expires > CURRENT_TIMESTAMP)
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
			s.log.Errorf("Failed to scan server user groups: %v\n", err)
			continue
		}

		if expires.Valid {
			group.Expires = expires.Time
		} else {
			group.Expires = time.Time{}
		}

		s.log.Debugf("Loading group '%s'[id: %d] for user: %d\n", group.GroupName, group.GroupID, UserID)
		groups = append(groups, group)
	}

	return groups, rows.Err()
}

func (s *Storage) loadUserPermissions(ctx context.Context, tx *sql.Tx, serverID int, UserID model.UserID) ([]model.UserPermission, error) {
	query := `
		SELECT
			permission,
			expires
		FROM server_user_permissions 
		WHERE steamid64 = ? AND server_id = ? AND (expires IS NULL OR expires = 0 OR expires > CURRENT_TIMESTAMP)
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
			s.log.Errorf("Failed to scan server user manager: %v\n", err)
			continue
		}

		if expires.Valid {
			permission.Expires = expires.Time
		} else {
			permission.Expires = time.Time{}
		}

		s.log.Debugf("Loading permission '%s' for user: %d\n", permission.Permission, UserID)
		permissions = append(permissions, permission)
	}

	s.log.Debugf("Loaded %d permissions\n", len(permissions))

	return permissions, rows.Err()
}

func (s *Storage) loadUserCookies(ctx context.Context, tx *sql.Tx, serverID int, UserID model.UserID) (map[string]string, error) {
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
			s.log.Errorf("Failed to scan server user cookies: %v\n", err)
			continue
		}

		s.log.Debugf("Loading cookie '%s' for user: %d\n", cookieKey, UserID)
		cookies[cookieKey] = cookieValue
	}

	s.log.Debugf("Loaded %d cookies\n", len(cookies))

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

func expiry(t time.Time) any {
	if t.IsZero() || t.Unix() <= 0 {
		return nil
	}
	return t.UTC()
}
