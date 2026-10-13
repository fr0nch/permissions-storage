# Changelog

## [0.1.1](https://github.com/fr0nch/permissions-storage/compare/v0.1.0...v0.1.1) (2026-10-13)


### Bug Fixes

* **ci:** build conda packages with rattler-build, allow uploading release assets ([ceeac30](https://github.com/fr0nch/permissions-storage/commit/ceeac3050dd88a695ee344ba3e87f164ead04c8e))

## [0.1.0](https://github.com/fr0nch/permissions-storage/compare/v0.0.1...v0.1.0) (2026-10-13)


### Features

* add initial plugin implementation ([9b7f964](https://github.com/fr0nch/permissions-storage/commit/9b7f964c1f59b1bcc3ccd52db6dedf68fcb213c7))
* add state and wildcard ([2fc835c](https://github.com/fr0nch/permissions-storage/commit/2fc835c8bda29adb29da826750f1a5b07c98ad44))
* implement database readiness checks and update storage structure ([23a8ca1](https://github.com/fr0nch/permissions-storage/commit/23a8ca17c0a0fc736f9a7ebfb6eb78924a947ebd))
* rework storage plugin ([94f54ff](https://github.com/fr0nch/permissions-storage/commit/94f54ff63e245ab971f0870fd1285121c3207058))
* update postgres and sqlite storages, protect group deletion by foreign keys ([1a36eef](https://github.com/fr0nch/permissions-storage/commit/1a36eef7e883cf3bef8c5fef3d571399d95f833b))


### Bug Fixes

* add plugin manifest ([b47f185](https://github.com/fr0nch/permissions-storage/commit/b47f185efec2240e48559340fb03b6cd60f6e41a))
* allow removing a group parent ([d23f46d](https://github.com/fr0nch/permissions-storage/commit/d23f46d9dc4d479bcffebcc7e00547bb6026ff88))
* comment out logging in config ([f0d0a8d](https://github.com/fr0nch/permissions-storage/commit/f0d0a8d762b19ce54ccb3fa07381a5a1fffe1f0b))
* group option change always being cancelled, return on DB not ready ([e6e1753](https://github.com/fr0nch/permissions-storage/commit/e6e1753fdf7dc0f6286e392e2eb7ebef4d893092))
* remove unused table_prefix setting ([855f94a](https://github.com/fr0nch/permissions-storage/commit/855f94ab6ab7d8f5c84a89d602d17e172caddc6d))
* sqlite always work as server 0, ignore server_id and global_cookie ([5b53327](https://github.com/fr0nch/permissions-storage/commit/5b533271415df0b22217472fc0e2e6d57a0a1549))
* sqlite timezone and permissions query ([e5a744c](https://github.com/fr0nch/permissions-storage/commit/e5a744cc604125797bc1d3711dd9a63a7ea3ef5c))
* update for new plugify version ([105ec3d](https://github.com/fr0nch/permissions-storage/commit/105ec3ddc1e4aaf2837a5867fe8133c5d41e63fc))
