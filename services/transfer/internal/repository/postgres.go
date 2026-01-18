package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	db "github.com/myacey/jxgercorp-banking/services/transfer/internal/repository/sqlc"
)

const ErrUniqueViolationCode = "23505"

type PostgresConfig struct {
	Host     string `mapstrucrute:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string // secret from env
	DBName   string `mapstructure:"db_name"`
}

func ConfiurePostgres(ctx context.Context, cfg PostgresConfig) (*db.Queries, *pgxpool.Pool, error) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	poolCfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("canot parse postgres config: %w", err)
	}
	poolCfg.MaxConns = 200
	poolCfg.MinConns = 50
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create postgres pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot ping postgres: %w", err)
	}

	queries := db.New(pool)

	return queries, pool, nil
}

func IsUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == ErrUniqueViolationCode
	}
	return false
}
