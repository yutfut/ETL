package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func NewPostgresPool(
	ctx context.Context,
	secret Postgres,
) (*pgxpool.Pool, error) {
	connectionUrl := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		secret.Host,
		secret.Port,
		secret.Username,
		secret.Password,
		secret.Database,
	)

	_cfg, err := pgxpool.ParseConfig(connectionUrl)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(
		ctx,
		_cfg,
	)
	if err != nil {
		log.Fatalf("unable to create connection pool: %s", err)
		return nil, err
	}

	err = pool.Ping(ctx)

	return pool, err
}
