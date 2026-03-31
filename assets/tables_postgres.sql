CREATE OR REPLACE FUNCTION update_updated_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_lastvisit_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.lastvisit_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    inheritance_id INTEGER DEFAULT NULL REFERENCES groups(id) ON DELETE SET NULL,
    CHECK (TRIM(name) <> '')
);

CREATE INDEX idx_groups_inheritance_id ON groups(inheritance_id);

CREATE TABLE IF NOT EXISTS group_options (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    option_key VARCHAR(255) NOT NULL,
    option_value TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, option_key),
    CHECK (TRIM(option_key) <> '')
);

CREATE INDEX idx_group_options_group_id ON group_options(group_id);
CREATE INDEX idx_group_options_option_key ON group_options(option_key);

CREATE TRIGGER update_group_options_updated
BEFORE UPDATE ON group_options
FOR EACH ROW
EXECUTE FUNCTION update_updated_column();

CREATE TABLE IF NOT EXISTS group_permissions (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission VARCHAR(255) NOT NULL,
    UNIQUE(group_id, permission),
    CHECK (TRIM(permission) <> '')
);

CREATE INDEX idx_group_permissions_group_id ON group_permissions(group_id);
CREATE INDEX idx_group_permissions_permission ON group_permissions(permission);

CREATE TABLE IF NOT EXISTS users (
    steamid64 BIGINT PRIMARY KEY,
    name VARCHAR(128) DEFAULT NULL,
    immunity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    lastvisit_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_users_lastvisit_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_lastvisit_at_column();

CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(32) DEFAULT NULL,
    default_group INTEGER DEFAULT NULL REFERENCES groups(id) ON DELETE SET NULL
);

CREATE INDEX idx_servers_default_group ON servers(default_group);

CREATE TABLE IF NOT EXISTS server_groups (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE NO ACTION,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE NO ACTION,
    UNIQUE(server_id, group_id)
);

CREATE INDEX idx_server_groups_server_id ON server_groups(server_id);
CREATE INDEX idx_server_groups_group_id ON server_groups(group_id);

INSERT INTO groups (id, name, priority, inheritance_id) VALUES (1, 'Default', 0, null);
INSERT INTO servers (id, name, address, default_group) VALUES (0, 'All Servers', null, 1);
INSERT INTO server_groups (id, server_id, group_id) VALUES (1, 0, 1);

CREATE TABLE IF NOT EXISTS server_user_groups (
    id SERIAL PRIMARY KEY,
    steamid64 BIGINT NOT NULL REFERENCES users(steamid64) ON DELETE CASCADE,
    server_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires TIMESTAMP NULL DEFAULT NULL,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(server_id, steamid64, group_id),
    FOREIGN KEY(server_id, group_id) REFERENCES server_groups(server_id, group_id) ON DELETE CASCADE
);

CREATE INDEX idx_server_user_groups_steamid64 ON server_user_groups(steamid64);
CREATE INDEX idx_server_user_groups_server_group ON server_user_groups(server_id, group_id);

CREATE TRIGGER update_server_user_groups_updated
BEFORE UPDATE ON server_user_groups
FOR EACH ROW
EXECUTE FUNCTION update_updated_column();

CREATE TABLE IF NOT EXISTS server_user_permissions (
    id SERIAL PRIMARY KEY,
    steamid64 BIGINT NOT NULL REFERENCES users(steamid64) ON DELETE CASCADE,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    permission VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires TIMESTAMP NULL DEFAULT NULL,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(steamid64, server_id, permission),
    CHECK (TRIM(permission) <> '')
);

CREATE INDEX idx_server_user_permissions_permission ON server_user_permissions(permission);
CREATE INDEX idx_server_user_permissions_server_id ON server_user_permissions(server_id);
CREATE INDEX idx_server_user_permissions_steamid64 ON server_user_permissions(steamid64);

CREATE TRIGGER update_server_user_permissions_updated
BEFORE UPDATE ON server_user_permissions
FOR EACH ROW
EXECUTE FUNCTION update_updated_column();

CREATE TABLE IF NOT EXISTS user_cookies (
    id SERIAL PRIMARY KEY,
    steamid64 BIGINT NOT NULL REFERENCES users(steamid64) ON DELETE CASCADE,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    option_key VARCHAR(255) NOT NULL,
    option_value TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(steamid64, server_id, option_key),
    CHECK (TRIM(option_key) <> '')
);

CREATE INDEX idx_user_cookies_servers ON user_cookies(server_id);
CREATE INDEX idx_user_cookies_steamid64 ON user_cookies(steamid64);
CREATE INDEX idx_user_cookies_option_key ON user_cookies(option_key);

CREATE TRIGGER update_user_cookies_updated
BEFORE UPDATE ON user_cookies
FOR EACH ROW
EXECUTE FUNCTION update_updated_column();