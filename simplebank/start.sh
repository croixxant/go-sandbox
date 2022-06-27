#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migrations -database "mysql://${DB_DSN}" -verbose up

echo "start the app"
exec "$@"