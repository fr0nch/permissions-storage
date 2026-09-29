package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"permissions-storage/internal/pkg/model"
	"slices"
	"strings"
	"time"
)

type Cookie struct {
	userID model.UserID
	key    string
	value  any
}

func (s *Storage) LoadUser(ctx context.Context, UserID model.UserID, username string) (*model.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `
		INSERT INTO users (steamid64, name, immunity, lastvisit_at)
		VALUES (?, ?, 0, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			lastvisit_at = VALUES(lastvisit_at)
		RETURNING immunity
	`

	user := new(model.User{UserID: UserID, Name: username})

	row := tx.QueryRowContext(ctx, query, UserID, username, time.Now())
	if err = row.Scan(&user.Immunity); err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
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

func (s *Storage) AddUserPermission(ctx context.Context, userID model.UserID, permission *model.UserPermission, wildcard bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	expires := expiry(permission.Expires)

	if wildcard {
		query := `UPDATE server_user_permissions
				  SET permission = ?, state = ?, expires = ?, updated = CURRENT_TIMESTAMP
				  WHERE steamid64 = ? AND server_id = ? AND permission = ?`

		searchPerm := strings.TrimSuffix(permission.Name, ".*")

		res, err := tx.ExecContext(ctx, query,
			permission.Name, permission.State, expires, userID, s.settings.ServerID, searchPerm,
		)

		if err != nil {
			return fmt.Errorf("failed to replace permission with wildcard: %w", err)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("check rows affected: %w", err)
		}

		if rows > 0 {
			return tx.Commit()
		}

		s.log.Warnf("Failed to update permission (not found): %s; inserting for user %v", searchPerm, userID)
	}

	query := `INSERT INTO server_user_permissions (steamid64, server_id, permission, state, expires)
			  VALUES (?, ?, ?, ?, ?)
			  ON DUPLICATE KEY UPDATE state = ?, expires = ?, updated = CURRENT_TIMESTAMP`

	_, err = tx.ExecContext(ctx, query,
		userID, s.settings.ServerID, permission.Name, permission.State, expires, permission.State, expires,
	)
	if err != nil {
		return fmt.Errorf("could not add permission: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) RemoveUserPermission(ctx context.Context, userID model.UserID, permission *model.UserPermission) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	var query string
	if s.settings.DeleteRemovedRecords {
		query = `DELETE FROM server_user_permissions WHERE steamid64 = ? AND server_id = ? AND permission = ?`
	} else {
		query = `UPDATE server_user_permissions
				 SET expires = CURRENT_TIMESTAMP, updated = CURRENT_TIMESTAMP
				 WHERE steamid64 = ? AND server_id = ? AND permission = ?`
	}

	_, err = tx.ExecContext(ctx, query, userID, s.settings.ServerID, permission.Name)
	if err != nil {
		return fmt.Errorf("could not remove/deactivate permission: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) AddUserGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO server_user_groups (steamid64, server_id, group_id, expires)
			  VALUES (?, ?, ?, ?)
			  ON DUPLICATE KEY UPDATE expires = ?`

	_, err = tx.ExecContext(ctx, query, userID, s.settings.ServerID, group.ID, expiry(group.Expires), expiry(group.Expires))
	if err != nil {
		return fmt.Errorf("could not add group: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) RemoveUserGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	var query string
	if s.settings.DeleteRemovedRecords {
		query = `DELETE FROM server_user_groups WHERE steamid64 = ? AND server_id = ? AND group_id = ?`
	} else {
		query = `UPDATE server_user_groups
				 SET expires = CURRENT_TIMESTAMP, updated = CURRENT_TIMESTAMP
				 WHERE steamid64 = ? AND server_id = ? AND group_id = ?`
	}

	_, err = tx.ExecContext(ctx, query, userID, s.settings.ServerID, group.ID)
	if err != nil {
		return fmt.Errorf("could not remove/deactivate group: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) SetCookie(ctx context.Context, userID model.UserID, key string, value any) {
	s.addQuery(ctx, Cookie{
		userID: userID,
		key:    key,
		value:  value,
	})
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
			expires
		FROM server_user_groups
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
		var groupID model.GroupID
		var expires sql.NullTime

		if err := rows.Scan(&groupID, &expires); err != nil {
			s.log.Errorf("Failed to scan server user groups: %v\n", err)
			continue
		}

		group.Group = new(model.Group{ID: groupID})

		if expires.Valid {
			group.Expires = expires.Time
		} else {
			group.Expires = time.Time{}
		}

		s.log.Debugf("Loading group '%s'[id: %d] for user: %d\n", group.Name, group.ID, UserID)
		groups = append(groups, group)
	}

	return groups, rows.Err()
}

func (s *Storage) loadUserPermissions(ctx context.Context, tx *sql.Tx, serverID int, UserID model.UserID) ([]model.UserPermission, error) {
	query := `
		SELECT
			permission,
			state,
			expires
		FROM server_user_permissions 
		WHERE steamid64 = ? AND server_id = ? AND (expires IS NULL OR expires = 0 OR expires > CURRENT_TIMESTAMP())
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

		if err := rows.Scan(&permission.Name, &permission.State, &expires); err != nil {
			s.log.Errorf("Failed to scan server user manager: %v\n", err)
			continue
		}

		if expires.Valid {
			permission.Expires = expires.Time
		} else {
			permission.Expires = time.Time{}
		}

		s.log.Debugf("Loading permission '%s'[state:%v] for user: %d\n", permission.Name, permission.State, UserID)
		permissions = append(permissions, permission)
	}

	s.log.Debugf("Loaded %d permissions\n", len(permissions))

	return permissions, rows.Err()
}

func (s *Storage) loadUserCookies(ctx context.Context, tx *sql.Tx, serverID int, UserID model.UserID) (model.Cookie, error) {
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

	var cookies = make(model.Cookie)

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

func (s *Storage) addQuery(ctx context.Context, cookie Cookie) {
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
			err := s.setCookie(ctx, _cookie.userID, _cookie.key, _cookie.value)
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
			  ON DUPLICATE KEY UPDATE option_value = ?`

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

func expiry(t time.Time) any {
	if t.IsZero() || t.Unix() <= 0 {
		return nil
	}
	return t.UTC()
}
