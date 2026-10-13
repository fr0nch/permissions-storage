package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"permissions-storage/internal/pkg/config"
	"sync"
	"time"

	"github.com/fr0nch/logger"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB

	ready  chan struct{}
	failed chan struct{}

	settings *config.Settings
	log      *logger.Logger

	delayedQueryTimer *time.Timer
	cookies           []Cookie
	closed            bool
	mu                sync.Mutex
}

func NewStorage(dsn url.URL, settings *config.Settings, logger *logger.Logger) (*Storage, error) {
	s := new(Storage{
		ready:    make(chan struct{}),
		failed:   make(chan struct{}),
		settings: settings,
		log:      logger,
	})

	dir := filepath.Dir(dsn.Path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	db, err := sql.Open("sqlite3", fileDSN(dsn))
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

func fileDSN(dsn url.URL) string {
	params := dsn.Query()

	setDefault := func(key, value string) {
		if !params.Has(key) {
			params.Set(key, value)
		}
	}

	setDefault("_foreign_keys", "on")
	setDefault("_busy_timeout", "5000")
	setDefault("_txlock", "immediate")

	u := url.URL{Scheme: "file", Path: dsn.Path, RawQuery: params.Encode()}

	return u.String()
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
	s.mu.Lock()
	s.closed = true
	if s.delayedQueryTimer != nil {
		s.delayedQueryTimer.Stop()
	}
	s.mu.Unlock()

	s.execQueries()
	s.db.Close()
}

func (s *Storage) CreateTables(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	return tx.Commit()
}
