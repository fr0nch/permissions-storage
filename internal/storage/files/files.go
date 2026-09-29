package files

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"permissions-storage/internal/pkg/config"
	"permissions-storage/internal/pkg/model"
	"strings"
	"sync"
	"time"

	"github.com/fr0nch/logger"
)

type Storage struct {
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
