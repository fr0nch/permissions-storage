package config

type Config struct {
	Settings *Settings `yaml:"settings"`
}

type ConfigsType interface {
	Settings
}

type PoolSettings struct {
	MaximumPoolSize     int `yaml:"maximum_pool_size" head_comment:"Sets the maximum size of the MySQL connection pool.\n- Basically this value will determine the maximum number of actual\n  connections to the database backend.\n- More information about determining the size of connection pools can be found here:\n  https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing"`
	MaximumIdle         int `yaml:"maximum_idle" head_comment:"Sets the maximum number of idle connections in the pool.\n- The pool will not keep more than this number of idle connections, but may keep fewer.\n- Setting this value equal to 'maximum_pool_size' makes the pool behave closer to a fixed-size pool and can improve performance under sudden load spikes."`
	MaximumLifetime     int `yaml:"maximum_lifetime" line_comment:"30 minutes" head_comment:"This setting controls the maximum lifetime of a connection in the pool in seconds.\n- The value should be at least 30 seconds less than any database or infrastructure imposed\n  connection time limit."`
	MaximumIdleLifetime int `yaml:"maximum_idle_lifetime" head_comment:"This setting controls the maximum idle lifetime of a connection in the pool in seconds.\n- Idle connections that exceed this duration will be closed and removed from the pool.\n- This helps to prevent resource waste by closing unused connections and ensures that stale\n  connections (which may have been dropped by the server or firewall) are not reused.\n- Recommended value: 2-5 minutes (120-300 seconds).\n- Set this value to 0 to disable idle connection timeout (not recommended for production)."`
}

type Database struct {
	Host string `yaml:"host" head_comment:"Define the host and port for the database.\n- The standard DB engine port is used by default\n  (MySQL: 3306, PostgreSQL: 5432)"`
	Port int    `yaml:"port"`

	Username string `yaml:"username" head_comment:"Credentials for the database."`
	Password string `yaml:"password"`

	PoolSettings PoolSettings `yaml:"pool_settings" head_comment:"These settings apply to the database connection pool.\n- The default values will be suitable for the majority of users.\n- Do not change these settings unless you know what you're doing!"`

	Database string `yaml:"database" head_comment:"The name of the database to store Plugify Permissions data in."`

	Schema string `yaml:"schema" head_comment:"PostgreSQL schema (namespace for tables)\n- Only used when using PostgreSQL storage. Ignored for MySQL and SQLite.\n- Default: public\n- Set this to separate plugin data inside one database"`

	Params map[string]string `yaml:"params" head_comment:"This setting allows you to define additional DSN parameters for the database connection.\n \nBy default for MySQL, the following options are applied automatically:\n - parseTime=true    : ensures proper parsing of DATE/TIME columns into time values\n - loc=Local         : sets the timezone for time values to the local system timezone\n \nOptional parameters you may want to include:\n - sslmode=disable   : for PostgreSQL, disables SSL (only if you understand the security implications)\n - _journal_mode=WAL : for SQLite, enables WAL mode for better concurrency\n - busy_timeout=5000 : for SQLite, sets a timeout for database locks (milliseconds)\n \nAdditional MySQL options you might set if needed:\n - collation=utf8mb4_general_ci : ensure full UTF-8 support including emojis\n - timeout=5s                  : set connection timeout\n \nExample usage:\nparams:\n sslmode: disable\n _journal_mode: WAL\n busy_timeout: 5000"`

	TablePrefix string `yaml:"table_prefix" head_comment:"The prefix for all Plugify Permissions SQL tables.\n \n- This only applies for remote SQL storage types (MySQL, MariaDB, etc).\n- Change this if you want to use different tables for different servers."`
}

type Settings struct {
	ServerID             int  `yaml:"server_id" head_comment:"Unique identifier of the server.\n- If set to 0, the plugin will operate on all available servers.\n  This allows for global operations without specifying a single server.\n- Note: This parameter is ignored when using SQLite storage,\n  since SQLite does not support multi-server separation in a single database."`
	GlobalCookie         bool `yaml:"global_cookie" head_comment:"Determines whether to use a global cookie shared across all servers, rather than individual server-specific cookies.\n- Useful for synchronizing user sessions or settings globally.\n- Note: This setting has no effect with SQLite storage."`
	DeleteExpiredRecords bool `yaml:"delete_expired_records" head_comment:"Automatically removes records from the database once the associated permissions or groups have expired.\n- Helps keep the database clean and ensures outdated entries do not persist."`
	DeleteRemovedRecords bool `yaml:"delete_removed_records" head_comment:"Automatically deletes database records when permissions or groups are manually removed from a player.\n- Ensures database consistency and prevents orphaned entries."`

	StorageMethod string   `yaml:"storage_method" head_comment:"How the plugin should store data\n \n- Possible options:\n \n  |  Remote databases - require connection information to be configured below\n  |=> MySQL/MariaDB (mysql)\n  |=> PostgreSQL (postgres)\n \n  |  Flatfile/local database - don't require any extra configuration\n  |=> SQLite (sqlite3)\n \n- A SQLite database is the default option."`
	Data          Database `yaml:"data" head_comment:"The following block defines the settings for remote database storage methods.\n \n- You don't need to touch any of the settings here if you're using a local storage method!\n- The connection detail options are shared between all remote storage types."`
}
