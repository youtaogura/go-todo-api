#!/bin/bash

# データベースが起動するまで待機
until nc -z -v -w30 $DB_HOST $DB_PORT
do
  echo "Waiting for database connection..."
  sleep 1
done

echo "Database is up - executing command"
go run gen/generator.go
exec air -c .air.toml