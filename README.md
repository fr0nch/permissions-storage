[![Русский](https://img.shields.io/badge/Русский-%F0%9F%87%B7%F0%9F%87%BA-green?style=for-the-badge)](README_ru.md)

# permissions-storage

A [Plugify](https://github.com/untrustedmodders/plugify) plugin that stores the data of the [permissions](https://github.com/untrustedmodders/plugify-plugin-permissions) plugin: groups, user permissions, user groups and cookies. Supported databases: MySQL/MariaDB, PostgreSQL and SQLite.

## Database setup

> [!IMPORTANT]
> **MySQL/MariaDB and PostgreSQL tables are not created by the plugin.** Create them yourself before the first start:
>
> ```sh
> mysql -u <user> -p <database> < assets/tables_mysql.sql
> psql -U <user> -d <database> -f assets/tables_postgres.sql
> ```
>
> SQLite needs no setup: the plugin creates the database file and its tables on start.

## Configuration

The config is created on the first start: `configs/permissions/storage.yml`.

| Setting | Description |
|---|---|
| `storage_method` | `mysql`, `postgres` or `sqlite3` |
| `data` | Connection to MySQL/PostgreSQL: host, port, credentials, database, pool settings. Not used by SQLite |
| `server_id` | Id of this server. `0` means all servers; groups and permissions given on server `0` apply on every server |
| `global_cookie` | Share cookies between all servers instead of keeping them per server |
| `delete_expired_records` | When a group or permission of a user expires, remove it in the database the same way as a manual removal (see `delete_removed_records`) |
| `delete_removed_records` | Delete removed groups and permissions instead of marking them expired |

The SQLite database is stored in `data/permissions/permissions.db`. SQLite always works as server `0`: `server_id` and `global_cookie` are ignored.

## Notes

- Removing a group or permission from a user removes it on all servers. To take a permission away on one server only, give it there with a deny (`-permission`).
- A group cannot be deleted while it has active users or is the default group of a server.

## Dependencies

- [permissions](https://github.com/untrustedmodders/plugify-plugin-permissions) `>= 1.3.0`

## Installation

From the conda channel:
```sh
conda install -c https://fr0nch.github.io/permissions-storage/ permissions-storage
```

## Building

Natively (requires Go and `jq`):
```sh
./build.sh [debug|release]
```

Via Docker (the image is built automatically on first run):
```sh
./build.sh [debug|release] --docker
```

Output in both cases: `build/<mode>/libpermissions-storage.so` and `build/<mode>/permissions-storage.pplugin`.

## License

Licensed under the [MIT License](LICENSE).
