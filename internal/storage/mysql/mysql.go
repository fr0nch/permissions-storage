package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"permissions-storage/internal/pkg/config"
	"strings"
	"sync"
	"time"

	"github.com/fr0nch/logger"
	_ "github.com/go-sql-driver/mysql"
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

func NewStorage(dsn url.URL, settings *config.Settings, logger *logger.Logger) (*Storage, error) {
	s := new(Storage{
		ready:    make(chan struct{}),
		failed:   make(chan struct{}),
		settings: settings,
		log:      logger,
	})

	db, err := sql.Open("mysql", strings.TrimPrefix(dsn.String(), "mysql://"))
	if err != nil {
		logger.Errorf("Doesn't open database: %v\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(settings.Data.PoolSettings.MaximumPoolSize)
	db.SetMaxIdleConns(settings.Data.PoolSettings.MaximumIdle)
	db.SetConnMaxIdleTime(time.Duration(settings.Data.PoolSettings.MaximumIdleLifetime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(settings.Data.PoolSettings.MaximumLifetime) * time.Second)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			db.Close()
			close(s.failed)

			logger.Errorf("Database is unreachable: %v", err)

			return
		}

		s.db = db

		_, err = db.Exec("SET time_zone = '+00:00'")
		if err != nil {
			logger.Errorf("Failed to set session timezone to UTC: %v\n", err)
			return
		}

		logger.Debugf("Database connection established successfully")
		close(s.ready)
	}()

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
	/*tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start the transaction.: %w", err)
	}

	defer tx.Rollback()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS groups (
	       id int(11) NOT NULL AUTO_INCREMENT,
	       name varchar(255) NOT NULL,
	       priority int(11) NOT NULL DEFAULT 0,
	       inheritance_id int(11) DEFAULT NULL,
	       PRIMARY KEY (id),
	       KEY fk_groups_groups (inheritance_id),
	       CONSTRAINT fk_groups_groups FOREIGN KEY (inheritance_id) REFERENCES groups (id) ON DELETE SET NULL ON UPDATE NO ACTION
	   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS group_options (
	       id int(11) NOT NULL AUTO_INCREMENT,
	       group_id int(11) NOT NULL,
	       option_key varchar(255) NOT NULL,
	       option_value longtext DEFAULT NULL,
	       created_at timestamp NOT NULL DEFAULT current_timestamp(),
	       updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	       PRIMARY KEY (id) USING BTREE,
	       UNIQUE KEY unique_group_option (group_id,option_key) USING BTREE,
	       KEY group_id (group_id),
	       KEY option_key (option_key),
	       CONSTRAINT fk_group_options_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE NO ACTION
	   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS group_permissions (
	       id int(11) NOT NULL AUTO_INCREMENT,
	       group_id int(11) NOT NULL,
	       permission varchar(255) NOT NULL,
	       PRIMARY KEY (id) USING BTREE,
	       UNIQUE KEY unique_group_permission (group_id,permission),
	       KEY group_id (group_id),
	       KEY permission (permission),
	       CONSTRAINT fk_group_permissions_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE NO ACTION
	   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS users (
	       ` + model.UserIDSQLColumnName + ` ` + model.UserIDSQLType + ` NOT NULL,
	       name varchar(128) DEFAULT NULL,
	       immunity int(11) NOT NULL DEFAULT 0,
	       created_at timestamp NOT NULL DEFAULT current_timestamp(),
	       lastvisit_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	       PRIMARY KEY (steamid64)
	   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS servers (
	       id int(11) NOT NULL AUTO_INCREMENT,
	       name varchar(255) NOT NULL,
	       address varchar(32) DEFAULT NULL,
	       default_group int(11) DEFAULT NULL,
	       PRIMARY KEY (id),
	       KEY fk_server_default_group (default_group),
	       CONSTRAINT fk_server_default_group FOREIGN KEY (default_group) REFERENCES groups (id) ON DELETE SET NULL ON UPDATE NO ACTION
	   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS server_groups (
	       id int(11) NOT NULL AUTO_INCREMENT,
	       server_id int(11) NOT NULL,
	       group_id int(11) NOT NULL,
	       PRIMARY KEY (id),
	       UNIQUE KEY unique_server_groups (server_id,group_id) USING BTREE,
	       KEY server_id (server_id),
	       KEY group_id (group_id) USING BTREE,
	       CONSTRAINT fk_server_groups_groups FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	       CONSTRAINT fk_server_groups_servers FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE NO ACTION ON UPDATE NO ACTION
	   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;`,

		`CREATE TABLE IF NOT EXISTS server_user_groups (
	       id int(11) NOT NULL AUTO_INCREMENT,
	       ` + model.UserIDSQLColumnName + ` ` + model.UserIDSQLType + ` NOT NULL,
	       server_id int(11) NOT NULL,
	       group_id int(11) NOT NULL,
	       created_at timestamp NOT NULL DEFAULT current_timestamp(),
	       expires timestamp NULL DEFAULT '0000-00-00 00:00:00',
	       updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	       PRIMARY KEY (id),
	       UNIQUE KEY unique_server_user_groups (server_id,` + model.UserIDSQLColumnName + `,group_id),
	       KEY fk_server_user_groups_` + model.UserIDSQLColumnName + ` (` + model.UserIDSQLColumnName + `),
	       KEY fk_server_user_groups_server_groups (server_id,group_id),
	       CONSTRAINT fk_server_user_groups_server_groups FOREIGN KEY (server_id, group_id) REFERENCES server_groups (server_id, group_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	       CONSTRAINT fk_server_user_groups_` + model.UserIDSQLColumnName + ` FOREIGN KEY (` + model.UserIDSQLColumnName + `) REFERENCES users (` + model.UserIDSQLColumnName + `) ON DELETE CASCADE ON UPDATE NO ACTION
	   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS server_user_permissions (
	       id int(11) NOT NULL AUTO_INCREMENT,
	       ` + model.UserIDSQLColumnName + ` ` + model.UserIDSQLType + ` NOT NULL,
	       server_id int(11) NOT NULL,
	       permission varchar(255) NOT NULL,
	       created_at timestamp NOT NULL DEFAULT current_timestamp(),
	       expires timestamp NULL DEFAULT '0000-00-00 00:00:00',
	       updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	       PRIMARY KEY (id),
	       UNIQUE KEY unique_server_user_permissions (` + model.UserIDSQLColumnName + `,server_id,permission) USING BTREE,
	       KEY permission (permission),
	       KEY server_id (server_id),
	       KEY ` + model.UserIDSQLColumnName + ` (` + model.UserIDSQLColumnName + `),
	       CONSTRAINT fk_user_permissions_server FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE ON UPDATE NO ACTION,
	       CONSTRAINT fk_user_permissions_` + model.UserIDSQLColumnName + ` FOREIGN KEY (` + model.UserIDSQLColumnName + `) REFERENCES users (` + model.UserIDSQLColumnName + `) ON DELETE CASCADE ON UPDATE NO ACTION,
	       CONSTRAINT chk_not_empty CHECK (trim(permission) <> '')
	   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS user_cookies (
	       id int(11) NOT NULL AUTO_INCREMENT,
	       ` + model.UserIDSQLColumnName + ` ` + model.UserIDSQLType + ` NOT NULL,
	       server_id int(11) NOT NULL,
	       option_key varchar(255) NOT NULL,
	       option_value longtext DEFAULT NULL,
	       created_at timestamp NOT NULL DEFAULT current_timestamp(),
	       updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	       PRIMARY KEY (id),
	       UNIQUE KEY unique_user_server_option (` + model.UserIDSQLColumnName + `,server_id,option_key),
	       KEY fk_user_cookies_servers (server_id),
	       KEY fk_user_cookies_` + model.UserIDSQLColumnName + ` (` + model.UserIDSQLColumnName + `),
	       KEY option_key (option_key),
	       CONSTRAINT fk_user_cookies_servers FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE ON UPDATE NO ACTION,
	       CONSTRAINT fk_user_cookies_` + model.UserIDSQLColumnName + ` FOREIGN KEY (` + model.UserIDSQLColumnName + `) REFERENCES users (` + model.UserIDSQLColumnName + `) ON DELETE CASCADE ON UPDATE NO ACTION
	   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	for _, query := range queries {
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("Query execution error:\n%s\n\nError: %w", query, err)
		}
	}

	return tx.Commit()*/
	return nil
}
