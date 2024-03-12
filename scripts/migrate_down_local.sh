export SHOPIFYX="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" && migrate -database ${SHOPIFYX} -path ./db/migrations down
