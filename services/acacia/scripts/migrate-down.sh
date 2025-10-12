export $(cat .env | xargs)

MIGRATION_DIR="./migrations"
NETWORK=acacia_dashboard-network

docker run -it -v $MIGRATION_DIR:/migrations --network $NETWORK migrate/migrate -path=/migrations/ -database $DATABASE_URL down $1
