package connections

import (
	"context"
	"fmt"
	"shopifyx/configs"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgConn() (*pgxpool.Pool, error) {
	ctx := context.Background()
	config := configs.LoadConfig()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?max_conns=-1", config.Database.User, config.Database.Pass, config.Database.Host, config.Database.Port, config.Database.Name)

	dbconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	dbconfig.MaxConnLifetime = 1 * time.Hour
	dbconfig.MaxConnIdleTime = 30 * time.Minute
	dbconfig.HealthCheckPeriod = 5 * time.Second
	dbconfig.MaxConns = -1
	dbconfig.MinConns = 5

	return pgxpool.NewWithConfig(ctx, dbconfig)
}
