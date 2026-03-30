CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    inheritance_id INTEGER DEFAULT NULL,
    FOREIGN KEY (inheritance_id) REFERENCES groups (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_groups_inheritance_id ON groups(inheritance_id);

CREATE TABLE IF NOT EXISTS group_options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    option_key TEXT NOT NULL,
    option_value TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, option_key),
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_group_options_group_id ON group_options(group_id);
CREATE INDEX IF NOT EXISTS idx_group_options_key ON group_options(option_key);

CREATE TABLE IF NOT EXISTS group_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    permission TEXT NOT NULL,
    UNIQUE(group_id, permission),
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_group_permissions_group_id ON group_permissions(group_id);
CREATE INDEX IF NOT EXISTS idx_group_permissions_permission ON group_permissions(permission);

CREATE TABLE IF NOT EXISTS users (
    steamid64 INTEGER PRIMARY KEY,
    name TEXT DEFAULT NULL,
    immunity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    lastvisit_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    address TEXT DEFAULT NULL,
    default_group INTEGER DEFAULT NULL,
    FOREIGN KEY (default_group) REFERENCES groups (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_servers_default_group ON servers(default_group);

CREATE TABLE IF NOT EXISTS server_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    UNIQUE(server_id, group_id),
    FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE NO ACTION,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_server_groups_server_id ON server_groups(server_id);
CREATE INDEX IF NOT EXISTS idx_server_groups_group_id ON server_groups(group_id);

CREATE TABLE IF NOT EXISTS server_user_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    steamid64 INTEGER NOT NULL,
    server_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires TIMESTAMP DEFAULT NULL,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(server_id, steamid64, group_id),
    FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE,
    FOREIGN KEY (server_id, group_id) REFERENCES server_groups (server_id, group_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_server_user_groups_steamid64 ON server_user_groups(steamid64);
CREATE INDEX IF NOT EXISTS idx_server_user_groups_server_group ON server_user_groups(server_id, group_id);

CREATE TABLE IF NOT EXISTS server_user_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    steamid64 INTEGER NOT NULL,
    server_id INTEGER NOT NULL,
    permission TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires TIMESTAMP DEFAULT NULL,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(steamid64, server_id, permission),
    FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE,
    FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE,
    CHECK (TRIM(permission) != '')
);

CREATE INDEX IF NOT EXISTS idx_server_user_permissions_permission ON server_user_permissions(permission);
CREATE INDEX IF NOT EXISTS idx_server_user_permissions_server_id ON server_user_permissions(server_id);
CREATE INDEX IF NOT EXISTS idx_server_user_permissions_steamid64 ON server_user_permissions(steamid64);

CREATE TABLE IF NOT EXISTS user_cookies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    steamid64 INTEGER NOT NULL,
    server_id INTEGER NOT NULL,
    option_key TEXT NOT NULL,
    option_value TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(steamid64, server_id, option_key),
    FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE,
    FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_cookies_server_id ON user_cookies(server_id);
CREATE INDEX IF NOT EXISTS idx_user_cookies_steamid64 ON user_cookies(steamid64);
CREATE INDEX IF NOT EXISTS idx_user_cookies_option_key ON user_cookies(option_key);