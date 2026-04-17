package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"permissions-storage/pkg/config"
	"permissions-storage/pkg/model"
	"permissions-storage/storage"
	"permissions-storage/storage/mysql"
	"permissions-storage/storage/postgres"
	"permissions-storage/storage/sqlite"
	"runtime/debug"

	"github.com/fr0nch/logger"
	"github.com/untrustedmodders/go-plugify"
)

type PermissionsStoragePlugin struct {
	pluginID uint64

	config  *config.Config
	storage storage.Storage

	log *logger.Logger

	groups         []*model.Group
	defaultGroupID int
}

//var Plugin *PermissionsStoragePlugin

func NewPermissionsStoragePlugin() *PermissionsStoragePlugin {
	return &PermissionsStoragePlugin{}
}

func init() {
	plugin := NewPermissionsStoragePlugin()

	plugify.OnPluginStart(plugin.OnPluginStart)
	plugify.OnPluginEnd(plugin.OnPluginEnd)
	plugify.OnPluginPanic(plugin.OnPluginPanic)
}

func (p *PermissionsStoragePlugin) OnPluginStart() {
	var err error

	p.log, err = logger.NewWithOptions(logger.Options{Folder: "permissions"})
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
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
				TablePrefix: "perms_",
			},
		},
	)

	if err != nil {
		panic("Failed to initialize config: " + err.Error())
	}

	dsn := BuildDSN(p.config.Settings)

	switch dsn.Scheme {
	case "mysql":
		p.storage, err = mysql.NewStorage(dsn, p.config.Settings, p.log)
		break
	case "postgres":
		p.storage, err = postgres.NewStorage(dsn, p.config.Settings, p.log)
		break
	case "file":
		p.storage, err = sqlite.NewStorage(dsn, p.config.Settings, p.log)
		break
	}

	if err != nil {
		panic("Failed to initialize storage: " + err.Error())
	}

	if dsn.Scheme == "file" {
		if err = p.storage.CreateTables(context.Background()); err != nil {
			panic("Failed to create tables: " + err.Error())
		}
	}

	p.RegisteringEvents()
}

func (p *PermissionsStoragePlugin) OnPluginEnd() {
	p.storage.Close()
}

func (p *PermissionsStoragePlugin) OnPluginPanic() []byte {
	return debug.Stack()
}

func BuildDSN(settings *config.Settings) url.URL {
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
		u.Scheme = "file"

		u.User = nil
		u.Path = filepath.Join(plugify.DataDir, "permissions", "permissions.db")
	default:
		panic("unsupported database type:" + settings.StorageMethod)
	}

	u.RawQuery = values.Encode()

	return u
}

func main() {}
