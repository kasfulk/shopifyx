export SHOPIFYX="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" && migrate -database ${SHOPIFYX} -path ./internal/database/migrations down
