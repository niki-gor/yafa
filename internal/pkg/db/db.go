package db

import (
	"fmt"
	"yafa/internal/pkg/env"

	"github.com/jackc/pgx"
)

func getConnString(env map[string]string) string {
	dsn := "postgres://%s:%s@localhost:5432/%s?%s"
	params := "sslmode=disable"
	return fmt.Sprintf(
		dsn,
		env["POSTGRES_USER"],
		env["POSTGRES_PASSWORD"],
		env["POSTGRES_DB"],
		params,
	)
}

func New() (*pgx.ConnPool, error) {
	env, err := env.GetRequired()
	if err != nil {
		return nil, err
	}

	connString := getConnString(env)

	connConf, err := pgx.ParseConnectionString(connString)
	if err != nil {
		return nil, err
	}

	return pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 50,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
}
