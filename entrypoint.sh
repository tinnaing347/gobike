#!/bin/bash
set -e

until curl -sL -I {$DB_ADDRESS}ping ; do
    >&2 echo "influx is unavailable - sleeping"
    sleep 1
done

>&2 echo "influx is up"
curl -XPOST {$DB_ADDRESS}query --data-urlencode 'q=CREATE DATABASE "bikedb"'
go run main.go

exec "$@"