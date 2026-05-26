package database

import (
	"database/sql"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/shubham071122/collab/internal/config"
)

func ConnectPostgres(cfg *config.Config) *sql.DB {

	connConfig, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	connConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	db, err := sql.Open("pgx", stdlib.RegisterConnConfig(connConfig))

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("postgres connected")

	return db
}
