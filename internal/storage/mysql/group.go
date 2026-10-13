package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"permissions-storage/internal/pkg/model"
	"permissions-storage/internal/storage"
	"strings"

	driver "github.com/go-sql-driver/mysql"
)

func (s *Storage) LoadGroups(ctx context.Context) (groups []*model.Group, defaultGroupID model.GroupID, err error) {
	tx, err := s.db.BeginTx(ctx, new(sql.TxOptions{ReadOnly: true}))
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

	ids := make([]model.GroupID, len(groups))
	for i := range groups {
		ids[i] = groups[i].ID
		groups[i].Options = make(model.Option)
	}

	var permsByGroup map[model.GroupID][]model.Permission
	var optsByGroup map[model.GroupID]model.Option

	permsByGroup, err = s.loadGroupPermissions(ctx, tx, ids)
	if err != nil {
		return nil, 0, err
	}

	optsByGroup, err = s.loadGroupOptions(ctx, tx, ids)
	if err != nil {
		return nil, 0, err
	}

	for i := range groups {
		gId := groups[i].ID
		groups[i].Permissions = permsByGroup[gId]
		groups[i].Options = optsByGroup[gId]
	}

	return groups, defaultID, tx.Commit()
}

func (s *Storage) CreateGroup(ctx context.Context, group *model.Group) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	var query = `INSERT INTO groups (name, priority, inheritance_id) 
				VALUES (?, ?, ?)`

	result, err := tx.ExecContext(ctx, query, group.Name, group.Priority, group.InheritanceID)
	if err != nil {
		return fmt.Errorf("could not add group: %w", err)
	}

	insId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("could not get last insert ID: %w", err)
	}

	group.ID = model.GroupID(insId)

	query = `INSERT INTO group_permissions (group_id, permission, state)
			 VALUES (?, ?, ?)
			 ON DUPLICATE KEY UPDATE state = ?`

	if group.Permissions != nil {
		for _, permission := range group.Permissions {
			_, err := tx.ExecContext(ctx, query, group.ID, permission.Name, permission.State, permission.State)
			if err != nil {
				return fmt.Errorf("could not add permission: %w", err)
			}
		}
	}

	{
		query = `INSERT INTO server_groups (server_id, group_id)
				VALUES (?, ?)`
		_, err := tx.ExecContext(ctx, query, s.settings.ServerID, group.ID)
		if err != nil {
			return fmt.Errorf("could not attach group to server: %w", err)
		}
	}

	return tx.Commit()
}

func (s *Storage) DeleteGroup(ctx context.Context, groupID model.GroupID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	var query = `DELETE FROM server_user_groups
				 WHERE group_id = ? AND NOT (expires IS NULL OR expires = 0 OR expires > CURRENT_TIMESTAMP())`

	if _, err = tx.ExecContext(ctx, query, groupID); err != nil {
		return fmt.Errorf("could not delete expired group users: %w", err)
	}

	query = `DELETE FROM groups WHERE id = ?`

	res, err := tx.ExecContext(ctx, query, groupID)
	if err != nil {
		return fmt.Errorf("could not delete group %d: %w", groupID, deleteGroupError(err))
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("group %d: %w", groupID, storage.ErrGroupNotFound)
	}

	return tx.Commit()
}

func deleteGroupError(err error) error {
	const rowIsReferenced = 1451 // ER_ROW_IS_REFERENCED_2

	var myErr *driver.MySQLError
	if !errors.As(err, &myErr) || myErr.Number != rowIsReferenced {
		return err
	}

	switch {
	case strings.Contains(myErr.Message, "`fk_server_user_groups_server_groups`"):
		return storage.ErrGroupHasUsers
	case strings.Contains(myErr.Message, "`fk_server_default_group`"):
		return storage.ErrGroupIsDefault
	}

	return err
}

func (s *Storage) AddGroupPermission(ctx context.Context, groupID model.GroupID, permission *model.Permission, wildcard bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	if wildcard {
		query := `UPDATE group_permissions
				  SET permission = ?, state = ?
				  WHERE group_id = ? AND permission = ?`

		searchPerm := strings.TrimSuffix(permission.Name, ".*")

		res, err := tx.ExecContext(ctx, query,
			permission.Name, permission.State, groupID, searchPerm,
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

		s.log.Warnf("Failed to update permission (not found): %s; inserting for group %d", searchPerm, groupID)
	}

	query := `INSERT INTO group_permissions (group_id, permission, state)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE 
			state = ?`

	_, err = tx.ExecContext(ctx, query,
		groupID, permission.Name, permission.State, permission.State,
	)
	if err != nil {
		return fmt.Errorf("could not add permission: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) RemoveGroupPermission(ctx context.Context, groupID model.GroupID, permission *model.Permission) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	var query = `DELETE FROM group_permissions WHERE group_id = ? AND permission = ?`

	_, err = tx.ExecContext(ctx, query, groupID, permission.Name)
	if err != nil {
		return fmt.Errorf("could not remove permission from group: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) SetParentGroup(ctx context.Context, groupID model.GroupID, inheritanceID model.GroupID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	var query = `UPDATE groups SET inheritance_id = ? WHERE id = ?`

	_, err = tx.ExecContext(ctx, query, inheritanceID, groupID)
	if err != nil {
		return fmt.Errorf("could not set parent for group: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) SetOption(ctx context.Context, groupID model.GroupID, key string, value any) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO group_options (group_id, option_key, option_value)
			  VALUES (?, ?, ?)
			  ON DUPLICATE KEY UPDATE option_value = ?`

	_, err = tx.ExecContext(ctx, query, groupID, key, storage.ToString(value), storage.ToString(value))
	if err != nil {
		return fmt.Errorf("could not setup option for group: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) loadGroups(ctx context.Context, tx *sql.Tx, serverID int) ([]*model.Group, model.GroupID, error) {
	query := `
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
	var defaultID model.GroupID

	for rows.Next() {
		g := new(model.Group{})
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

func (s *Storage) loadGroupPermissions(ctx context.Context, tx *sql.Tx, groupIDs []model.GroupID) (map[model.GroupID][]model.Permission, error) {
	in, args := makeInClause(groupIDs)

	query := fmt.Sprintf(`
		SELECT group_id, permission, state
		FROM group_permissions
		WHERE group_id IN (%s)
		ORDER BY group_id;
	`, in)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query manager: %w", err)
	}

	defer rows.Close()

	out := make(map[model.GroupID][]model.Permission, len(groupIDs))
	for rows.Next() {
		var gid model.GroupID
		var permission model.Permission

		if err := rows.Scan(&gid, &permission.Name, &permission.State); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}

		out[gid] = append(out[gid], permission)
	}

	return out, rows.Err()
}

func (s *Storage) loadGroupOptions(ctx context.Context, tx *sql.Tx, groupIDs []model.GroupID) (map[model.GroupID]model.Option, error) {
	in, args := makeInClause(groupIDs)

	query := fmt.Sprintf(`
		SELECT group_id, option_key, COALESCE(option_value, '')
		FROM group_options
		WHERE group_id IN (%s)
		ORDER BY group_id;
	`, in)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query options: %w", err)
	}
	defer rows.Close()

	out := make(map[model.GroupID]model.Option, len(groupIDs))
	for rows.Next() {
		var gid model.GroupID
		var k, v string

		if err := rows.Scan(&gid, &k, &v); err != nil {
			return nil, fmt.Errorf("scan option: %w", err)
		}

		if out[gid] == nil {
			out[gid] = make(model.Option)
		}

		out[gid][k] = v
	}

	return out, rows.Err()
}

func makeInClause[T model.ID](ids []T) (string, []any) {
	parts := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))

	for _, id := range ids {
		parts = append(parts, "?")
		args = append(args, id)
	}

	return strings.Join(parts, ", "), args
}
