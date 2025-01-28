package postgresrepo

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	db "github.com/myacey/jxgercorp-banking/db/sqlc"
	"github.com/myacey/jxgercorp-banking/shared/backconfig"
)

func ConfiurePostgres(config backconfig.Config) (*db.Queries, *sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.PostgresHost, config.PostgresPort, config.PostgresUser,
		config.DBPassword, config.PostgresDBName)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("canot open postgres conn:", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot ping postgres:", err)
	}

	queries := db.New(conn)

	return queries, conn, nil
}

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
