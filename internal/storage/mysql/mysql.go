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

	params := dsn.Query()
	if !params.Has("time_zone") {
		params.Set("time_zone", "'+00:00'")
	}
	dsn.RawQuery = params.Encode()

	db, err := sql.Open("mysql", strings.TrimPrefix(dsn.String(), "mysql://"))
	if err != nil {
		logger.Errorf("Doesn't open database: %v\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(settings.Data.PoolSettings.MaximumPoolSize)
	db.SetMaxIdleConns(settings.Data.PoolSettings.MaximumIdle)
	db.SetConnMaxIdleTime(time.Duration(settings.Data.PoolSettings.MaximumIdleLifetime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(settings.Data.PoolSettings.MaximumLifetime) * time.Second)

	s.db = db

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			logger.Errorf("Database is unreachable: %v", err)

			db.Close()
			close(s.failed)

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
	return nil
}
