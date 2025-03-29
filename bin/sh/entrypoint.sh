#!/bin/sh
set -e

if [ "$1" = "api" ]; then
    migrate -path=./migrations -database=${GREENLIGHT_DB_DSN} up
    exec /app/bin/api
else
    exec "$@"
fi