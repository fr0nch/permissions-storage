[![English](https://img.shields.io/badge/English-%F0%9F%87%AC%F0%9F%87%A7-blue?style=for-the-badge)](README.md)

# permissions-storage

Плагин [Plugify](https://github.com/untrustedmodders/plugify), который хранит данные плагина [permissions](https://github.com/untrustedmodders/plugify-plugin-permissions): группы, права и группы игроков, куки. Поддерживаемые базы данных: MySQL/MariaDB, PostgreSQL и SQLite.

## Подготовка базы данных

> [!IMPORTANT]
> **Таблицы в MySQL/MariaDB и PostgreSQL плагин не создаёт.** Создайте их самостоятельно до первого запуска:
>
> ```sh
> mysql -u <user> -p <database> < assets/tables_mysql.sql
> psql -U <user> -d <database> -f assets/tables_postgres.sql
> ```
>
> Для SQLite ничего делать не нужно: плагин сам создаёт файл базы и таблицы при запуске.

## Настройка

Конфиг создаётся при первом запуске: `configs/permissions/storage.yml`.

| Параметр | Описание |
|---|---|
| `storage_method` | `mysql`, `postgres` или `sqlite3` |
| `data` | Подключение к MySQL/PostgreSQL: адрес, порт, учётные данные, база, настройки пула. Для SQLite не используется |
| `server_id` | Идентификатор сервера. `0` — все серверы: группы и права, выданные на сервере `0`, действуют на каждом сервере |
| `global_cookie` | Общие куки для всех серверов вместо отдельных для каждого |
| `delete_expired_records` | Когда у игрока истекает группа или право, снимать их в базе так же, как при ручном снятии (см. `delete_removed_records`) |
| `delete_removed_records` | Удалять снятые группы и права, а не помечать их истёкшими |

База SQLite хранится в `data/permissions/permissions.db`. SQLite всегда работает как сервер `0`: `server_id` и `global_cookie` не учитываются.

## Особенности

- Снятие группы или права с игрока снимает их на всех серверах. Чтобы забрать право только на одном сервере, выдайте его там с запретом (`-permission`).
- Группу нельзя удалить, пока в ней есть активные игроки или она назначена группой по умолчанию на каком-либо сервере.

## Зависимости

- [permissions](https://github.com/untrustedmodders/plugify-plugin-permissions) `>= 1.3.0`

## Установка

Из канала conda:
```sh
conda install -c https://fr0nch.github.io/permissions-storage/ permissions-storage
```

## Сборка

Нативно (нужен установленный Go и `jq`):
```sh
./build.sh [debug|release]
```

Через Docker (образ собирается автоматически при первом запуске):
```sh
./build.sh [debug|release] --docker
```

Результат в обоих случаях: `build/<mode>/libpermissions-storage.so` и `build/<mode>/permissions-storage.pplugin`.

## Лицензия

Распространяется по лицензии [MIT](LICENSE).
