#!/bin/sh

echo "Waiting for Postgres to start..."
./wait-for-it.sh db:5432

echo "Migrating the database"
tern migrate --migrations ./migrations

echo "Starting the server..."
./bitlygo
