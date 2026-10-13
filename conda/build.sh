#!/bin/bash
# !build.sh - For Linux builds

set -ex

mkdir -p $PREFIX/bin

cp bin/libpermissions-storage.so $PREFIX/bin/
cp permissions-storage.pplugin $PREFIX/

chmod 755 $PREFIX/bin/libpermissions-storage.so
chmod 644 $PREFIX/permissions-storage.pplugin
