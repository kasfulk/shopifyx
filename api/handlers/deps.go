package handlers

import (
	"shopifyx/configs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Cfg    configs.Config
	DbPool *pgxpool.Pool
}
