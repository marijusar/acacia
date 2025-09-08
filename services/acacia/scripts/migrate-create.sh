#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: ./migrate-create.sh <migration_name>"
    echo "Example: ./migrate-create.sh create_users_table"
    exit 1
fi

MIGRATION_NAME=$1
MIGRATION_DIR="./migrations"

docker run -v $PWD/$MIGRATION_DIR:/migrations migrate/migrate create -ext sql -dir /migrations -seq $MIGRATION_NAME

# Fix ownership and permissions of created files
sudo chown $USER:$USER $MIGRATION_DIR/*.sql
chmod 644 $MIGRATION_DIR/*.sql