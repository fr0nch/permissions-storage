package internal

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"permissions-storage/internal/pkg/config"
	"permissions-storage/internal/pkg/model"
	"permissions-storage/internal/storage"
	"permissions-storage/internal/storage/mysql"
	"permissions-storage/internal/storage/postgres"
	"permissions-storage/internal/storage/sqlite"

	"github.com/fr0nch/logger"
	"github.com/untrustedmodders/go-plugify"
)

type PermissionsStoragePlugin struct {
	Plugin   plugify.Plugin
	PluginID int64

	config  *config.Config
	storage storage.Storage

	log *logger.Logger

	groups         []*model.Group
	defaultGroupID model.GroupID
}

func NewPermissionsStoragePlugin() *PermissionsStoragePlugin {
	return &PermissionsStoragePlugin{}
}

func (p *PermissionsStoragePlugin) OnPluginStart() error {
	p.PluginID = int64(p.Plugin.ID())

	var err error
	p.log, err = logger.NewWithOptions(logger.Options{PluginName: p.Plugin.Name(), Folder: "permissions"})
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	p.config = config.NewConfig()
	p.config.Settings, err = config.LoadConfig("permissions", "storage.yml",
		config.Settings{
			StorageMethod: "sqlite3",
			Data: config.Database{
				Port: 3306,
				PoolSettings: config.PoolSettings{
					MaximumPoolSize:     25,
					MaximumIdle:         5,
					MaximumLifetime:     1800,
					MaximumIdleLifetime: 300,
				},
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	//fmt.Println(p.config.Settings)

	dsn, err := BuildDSN(p.config.Settings)
	if err != nil {
		return fmt.Errorf("failed to initialize database connection: %w", err)
	}

	switch dsn.Scheme {
	case "mysql":
		p.storage, err = mysql.NewStorage(dsn, p.config.Settings, p.log)
		break
	case "postgres":
		p.storage, err = postgres.NewStorage(dsn, p.config.Settings, p.log)
		break
	case "sqlite3":
		p.storage, err = sqlite.NewStorage(dsn, p.config.Settings, p.log)
		break
	}

	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	if err = p.storage.CreateTables(context.Background()); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	p.RegisteringEvents()
	return nil
}

func (p *PermissionsStoragePlugin) OnPluginUpdate(dt float32) error {
	return nil
}

func (p *PermissionsStoragePlugin) OnPluginEnd() error {
	p.storage.Close()
	return nil
}

func BuildDSN(settings *config.Settings) (url.URL, error) {
	values := url.Values{}
	for k, v := range settings.Data.Params {
		values.Add(k, v)
	}

	u := url.URL{
		User: url.UserPassword(settings.Data.Username, settings.Data.Password),
	}

	switch settings.StorageMethod {
	case "mysql":
		u.Scheme = "mysql"
		u.Host = fmt.Sprintf("tcp(%s:%d)", settings.Data.Host, settings.Data.Port)
		u.Path = settings.Data.Database

		if _, ok := settings.Data.Params["loc"]; !ok {
			values.Add("loc", "UTC")
		}

		if _, ok := settings.Data.Params["parseTime"]; !ok {
			values.Add("parseTime", "true")
		}

	case "postgres":
		u.Scheme = "postgres"
		u.Host = fmt.Sprintf("%s:%d", settings.Data.Host, settings.Data.Port)
		u.Path = settings.Data.Database

		if schema := settings.Data.Schema; schema != "" {
			values.Set("search_path", schema)
		}

	case "sqlite3":
		u.Scheme = "sqlite3"

		u.User = nil
		u.Path = filepath.Join(plugify.DataDir(), "permissions", "permissions.db")
	default:
		return url.URL{}, fmt.Errorf("unsupported database type: %s", settings.StorageMethod)
	}

	u.RawQuery = values.Encode()

	return u, nil
}
