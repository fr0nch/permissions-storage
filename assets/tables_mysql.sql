CREATE TABLE IF NOT EXISTS perms_groups (
    id int(11) NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    priority int(11) NOT NULL DEFAULT 0,
    inheritance_id int(11) DEFAULT NULL,
    PRIMARY KEY (id),
    KEY fk_groups_groups (inheritance_id),
    CONSTRAINT fk_groups_groups FOREIGN KEY (inheritance_id) REFERENCES groups (id) ON DELETE SET NULL ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS perms_group_options (
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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS perms_group_permissions (
    id int(11) NOT NULL AUTO_INCREMENT,
    group_id int(11) NOT NULL,
    permission varchar(255) NOT NULL,
    PRIMARY KEY (id) USING BTREE,
    UNIQUE KEY unique_group_permission (group_id,permission),
    KEY group_id (group_id),
    KEY permission (permission),
    CONSTRAINT fk_group_permissions_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS perms_users (
    steamid64 bigint(20) NOT NULL,
    name varchar(128) DEFAULT NULL,
    immunity int(11) NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT current_timestamp(),
    lastvisit_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (steamid64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS perms_servers (
    id int(11) NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    address varchar(32) DEFAULT NULL,
    default_group int(11) DEFAULT NULL,
    PRIMARY KEY (id),
    KEY fk_server_default_group (default_group),
    CONSTRAINT fk_server_default_group FOREIGN KEY (default_group) REFERENCES groups (id) ON DELETE SET NULL ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS perms_server_groups (
    id int(11) NOT NULL AUTO_INCREMENT,
    server_id int(11) NOT NULL,
    group_id int(11) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY unique_server_groups (server_id,group_id) USING BTREE,
    KEY server_id (server_id),
    KEY group_id (group_id) USING BTREE,
    CONSTRAINT fk_server_groups_groups FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT fk_server_groups_servers FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;

CREATE TABLE IF NOT EXISTS perms_server_user_groups (
    id int(11) NOT NULL AUTO_INCREMENT,
    steamid64 bigint(20) NOT NULL,
    server_id int(11) NOT NULL,
    group_id int(11) NOT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp(),
    expires timestamp NULL DEFAULT NULL,
    updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (id),
    UNIQUE KEY unique_server_user_groups (server_id,steamid64,group_id),
    KEY fk_server_user_groups_steamid64 (steamid64),
    KEY fk_server_user_groups_server_groups (server_id,group_id),
    CONSTRAINT fk_server_user_groups_server_groups FOREIGN KEY (server_id, group_id) REFERENCES server_groups (server_id, group_id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT fk_server_user_groups_steamid64 FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS perms_server_user_permissions (
    id int(11) NOT NULL AUTO_INCREMENT,
    steamid64 bigint(20) NOT NULL,
    server_id int(11) NOT NULL,
    permission varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp(),
    expires timestamp NULL DEFAULT NULL,
    updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (id),
    UNIQUE KEY unique_server_user_permissions (steamid64,server_id,permission) USING BTREE,
    KEY permission (permission),
    KEY server_id (server_id),
    KEY steamid64 (steamid64),
    CONSTRAINT fk_user_permissions_server FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT fk_user_permissions_steamid64 FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT chk_not_empty CHECK (trim(permission) <> '')
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS perms_user_cookies (
    id int(11) NOT NULL AUTO_INCREMENT,
    steamid64 bigint(20) NOT NULL,
    server_id int(11) NOT NULL,
    option_key varchar(255) NOT NULL,
    option_value longtext DEFAULT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp(),
    updated timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (id),
    UNIQUE KEY unique_user_server_option (steamid64,server_id,option_key),
    KEY fk_user_cookies_servers (server_id),
    KEY i_user_cookies_steamid64 (steamid64),
    KEY option_key (option_key),
    CONSTRAINT fk_user_cookies_servers FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT fk_user_cookies_steamid64 FOREIGN KEY (steamid64) REFERENCES users (steamid64) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;