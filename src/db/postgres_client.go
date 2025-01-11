package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresDb struct {
	dbpool *pgxpool.Pool
}

func NewPostgresDb(databaseUrl string) *PostgresDb {

	log.Info().Msgf("Create postgresql connection to %s", databaseUrl)

	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil || len(databaseUrl) == 0 {
		panic(fmt.Errorf("Cannot connect ot Postgres: %s, err: %v", databaseUrl, err))
	}

	return &PostgresDb{
		dbpool: pool,
	}
}

func (client *PostgresDb) Close() {
	client.dbpool.Close()
}
