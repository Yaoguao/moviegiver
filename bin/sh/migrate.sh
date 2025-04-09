#!/bin/sh

if [ "$1" = "down" ]; then
  shift
  migrate -path=/app/migrations -database=$MOVIEGIVER_DB_DSN down "${1:-1}"  # По умолчанию 1
else
  migrate -path=/app/migrations -database=$MOVIEGIVER_DB_DSN up
fi