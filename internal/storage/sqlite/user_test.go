package sqlite

import (
	"context"
	"net/url"
	"path/filepath"
	"permissions-storage/internal/pkg/config"
	"permissions-storage/internal/pkg/model"
	"testing"
	"time"

	"github.com/fr0nch/logger"
)

func openTestStorage(t *testing.T) *Storage {
	t.Helper()

	log, err := logger.NewWithOptions(logger.Options{
		PluginName: "permissions-storage-test",
		Folder:     t.TempDir(),
	})
	if err != nil {
		t.Fatalf("init test logger: %v", err)
	}

	settings := &config.Settings{
		Data: config.Database{
			PoolSettings: config.PoolSettings{MaximumPoolSize: 2},
		},
	}
	dsn := url.URL{
		Scheme: "sqlite3",
		Path:   filepath.Join(t.TempDir(), "permissions.db"),
	}

	store, err := NewStorage(dsn, settings, log)
	if err != nil {
		t.Fatalf("open sqlite storage: %v", err)
	}
	t.Cleanup(store.Close)

	if err := store.CreateTables(context.Background()); err != nil {
		t.Fatalf("create tables: %v", err)
	}
	return store
}

func readUserRow(t *testing.T, store *Storage, steamID model.UserID) (string, int32, time.Time) {
	t.Helper()

	var name string
	var immunity int32
	var lastVisit time.Time
	err := store.db.QueryRow(
		`SELECT name, immunity, lastvisit_at FROM users WHERE steamid64 = ?`,
		steamID,
	).Scan(&name, &immunity, &lastVisit)
	if err != nil {
		t.Fatalf("query user %d: %v", steamID, err)
	}
	return name, immunity, lastVisit.UTC()
}

func TestLoadUserPreservesStoredProfileWhenOffline(t *testing.T) {
	ctx := context.Background()
	store := openTestStorage(t)

	// Mirror PostgreSQL's BEFORE UPDATE lastvisit_at trigger so any spurious UPDATE is caught.
	if _, err := store.db.ExecContext(ctx, `
		CREATE TABLE user_updates (steamid64 INTEGER NOT NULL);
		CREATE TRIGGER users_touch_lastvisit AFTER UPDATE ON users FOR EACH ROW
		BEGIN
			INSERT INTO user_updates (steamid64) VALUES (NEW.steamid64);
		END;
	`); err != nil {
		t.Fatalf("create update audit trigger: %v", err)
	}

	existingID := model.UserID(76561198000000001)
	historicalVisit := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	if _, err := store.db.ExecContext(
		ctx,
		`INSERT INTO users (steamid64, name, immunity, lastvisit_at) VALUES (?, ?, ?, ?)`,
		existingID, "Alice", 50, historicalVisit,
	); err != nil {
		t.Fatalf("seed existing user: %v", err)
	}

	offlineUser, err := store.LoadUser(ctx, existingID, "", true)
	if err != nil {
		t.Fatalf("offline LoadUser existing user: %v", err)
	}
	if offlineUser.Name != "Alice" || offlineUser.Immunity != 50 {
		t.Fatalf("expected stored profile Alice/50, got %q/%d", offlineUser.Name, offlineUser.Immunity)
	}

	var updateCount int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM user_updates WHERE steamid64 = ?`, existingID).Scan(&updateCount); err != nil {
		t.Fatalf("query user_updates: %v", err)
	}
	if updateCount != 0 {
		t.Fatalf("expected 0 row updates on offline lookup, got %d", updateCount)
	}

	name, immunity, lastVisit := readUserRow(t, store, existingID)
	if name != "Alice" || immunity != 50 || !lastVisit.Equal(historicalVisit) {
		t.Fatalf("offline lookup mutated profile: got name=%q immunity=%d lastvisit=%v", name, immunity, lastVisit)
	}

	onlineUser, err := store.LoadUser(ctx, existingID, "AliceUpdated", false)
	if err != nil {
		t.Fatalf("online LoadUser existing user: %v", err)
	}
	if onlineUser.Immunity != 50 {
		t.Fatalf("expected online LoadUser immunity 50, got %d", onlineUser.Immunity)
	}

	name, immunity, lastVisit = readUserRow(t, store, existingID)
	if name != "AliceUpdated" || immunity != 50 || !lastVisit.After(historicalVisit) {
		t.Fatalf("online load did not update profile: got name=%q immunity=%d lastvisit=%v", name, immunity, lastVisit)
	}

	unseenID := model.UserID(76561198000000002)
	newOfflineUser, err := store.LoadUser(ctx, unseenID, "NewOffline", true)
	if err != nil {
		t.Fatalf("offline LoadUser unseen user: %v", err)
	}
	if newOfflineUser.Name != "NewOffline" || newOfflineUser.Immunity != 0 {
		t.Fatalf("expected unseen offline user NewOffline/0, got %q/%d", newOfflineUser.Name, newOfflineUser.Immunity)
	}

	name, immunity, _ = readUserRow(t, store, unseenID)
	if name != "NewOffline" || immunity != 0 {
		t.Fatalf("unexpected unseen offline row: got name=%q immunity=%d", name, immunity)
	}
}
