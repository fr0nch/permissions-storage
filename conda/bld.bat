@echo off
REM bld.bat - For Windows builds

if not exist "%PREFIX%\bin" mkdir "%PREFIX%\bin"

copy bin\libpermissions-storage.dll "%PREFIX%\bin\" || exit 1
copy permissions-storage.pplugin "%PREFIX%\" || exit 1

exit 0
