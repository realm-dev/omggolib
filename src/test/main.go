package main

import (
	"github.com/realm-dev/omggolib/src/db"
	"github.com/realm-dev/omggolib/src/model"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	DATABASE_URL := os.Getenv("DATABASE_URL")

	postgres := db.NewPostgresDb(DATABASE_URL)
	defer postgres.Close()

	newRequests, err := postgres.GetWithdrawalRequests(model.Requested)
	if err != nil {
		log.Panic().Msgf("Cannot get withdrawal requests: %v", err)
		return
	}

	for _, request := range newRequests {
		log.Info().Msgf("accountId: %d, lamports: %d", request.AccountId, request.Lamports)
	}
}
