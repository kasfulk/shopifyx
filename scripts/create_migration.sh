# sh scripts/create_migration.sh <migration_name>
# sh scripts/create_migration.sh create_table_users
migrate create -ext sql -dir ./db/migrations -seq $*