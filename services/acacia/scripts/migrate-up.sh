export $(cat .env | xargs)

MIGRATION_DIR="$PWD/migrations"
NETWORK=acacia_dashboard-network

docker run -v $MIGRATION_DIR:/migrations --network $NETWORK migrate/migrate -path=/migrations/ -database $DATABASE_URL up
