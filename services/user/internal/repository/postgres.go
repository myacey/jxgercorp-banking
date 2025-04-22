package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	db "github.com/myacey/jxgercorp-banking/services/user/internal/repository/sqlc"
)

const ErrUniqueViolationCode = "23505"

type PostgresConfig struct {
	Host     string `mapstrucrute:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string // secret from env
	DBName   string `mapstructure:"db_name"`
}

func ConfiurePostgres(cfg PostgresConfig) (*db.Queries, *sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User,
		cfg.Password, cfg.DBName)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("canot open postgres conn: %w", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot ping postgres: %w", err)
	}

	conn.SetMaxOpenConns(200)
	conn.SetMaxIdleConns(50)
	conn.SetConnMaxIdleTime(5 * time.Minute)

	queries := db.New(conn)

	return queries, conn, nil
}

func IsUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == ErrUniqueViolationCode
	}
	return false
}
