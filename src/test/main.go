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

	commissions, err := postgres.GetPaybackCommissions(99999, model.CS_Paid)
	log.Info().Msgf("found %d transactions", len(commissions))
	if err != nil {
		log.Panic().Msgf("Cannot cannot get payback commissions for account: %d, error: %v", 99999, err)
		return
	}

	for _, commission := range commissions {
		err = postgres.UpdateCommissionStatus(commission.Hash, model.CS_Paidback)
		if err != nil {
			log.Panic().Msgf("Cannot update commission for account: %d, error: %v", 99999, err)
			return
		}

		updatedCommission, err := postgres.GetCommission(commission.Hash)
		if err != nil || updatedCommission == nil {
			log.Panic().Msgf("failed to get commission status: %v", err)
			return
		}

		log.Info().Msgf("new commission status %d", updatedCommission.Status)
		err = postgres.UpdateCommissionStatus(commission.Hash, model.CS_Paid)
		if err != nil {
			log.Panic().Msgf("Cannot cannot get payback commissions for account: %d, error: %v", 99999, err)
			return
		}
	}

	newRequests, err := postgres.GetWithdrawalRequests(model.WS_Requested)
	if err != nil {
		log.Panic().Msgf("Cannot get withdrawal requests: %v", err)
		return
	}

	for _, request := range newRequests {
		log.Info().Msgf("accountId: %d, status: %d", request.AccountId, request.Status)

		err = postgres.SetWithdrawalCalculatedLamports(request.AccountId, 100)
		if err != nil {
			log.Panic().Msgf("Cannot set calculated lamports to withdrawal request account: %d, error: %v", request.AccountId, err)
			return
		}

		preparedRequest, err := postgres.GetWithdrawalRequest(request.AccountId, model.WS_Requested)
		if err != nil {
			log.Panic().Msgf("Cannot get withdrawal request for account: %d, error: %v", preparedRequest.AccountId, err)
			return
		}
		log.Info().Msgf("calculated %d lamports for accountId: %d", preparedRequest.Lamports, preparedRequest.AccountId)

		err = postgres.SetWithdrawalResult(request.AccountId, "213234j2k3kj43543u3")
		if err != nil {
			log.Panic().Msgf("Cannot set paid lamports to withdrawal request account: %d, error: %v", request.AccountId, err)
			return
		}

		completedRequest, err := postgres.GetWithdrawalRequest(request.AccountId, model.WS_PaidOut)
		if err != nil {
			log.Panic().Msgf("Cannot get withdrawal request for account: %d, error: %v", completedRequest.AccountId, err)
			return
		}
		log.Info().Msgf("withdrawal request completed for accountId: %d, paid lamports: %d", completedRequest.AccountId, completedRequest.Lamports)

	}

}
