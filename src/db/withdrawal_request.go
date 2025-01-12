package db

import (
	"context"
	"fmt"
	"github.com/realm-dev/omggolib/src/model"
	"os"
)

func (client *PostgresDb) InsertWithdrawalRequest(withdrawalRequest model.WithdrawalRequest) error {
	commandTag, err := client.dbpool.Exec(context.Background(), "INSERT INTO withdrawal_requests ("+
		"account_id, wallet_public_key, timestamp, status) "+
		"VALUES ($1, $2, $3, $4)",
		withdrawalRequest.AccountId, withdrawalRequest.WalletPublicKey, withdrawalRequest.Timestamp, withdrawalRequest.Status)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec failed: %v\n", err)
		return err
	}

	if commandTag.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
		return err
	}
	return nil
}

func (client *PostgresDb) SetNewWithdrawalRequestTimestamp(accountId int64, status model.WithdrawalStatus, timestamp int64) error {
	_, err := client.dbpool.Exec(context.Background(),
		"UPDATE withdrawal_requests SET timestamp = $1 WHERE account_id = $2 and status = $3",
		timestamp, accountId, int(status))

	return err
}

func (client *PostgresDb) GetWithdrawalRequest(accountId int64, status model.WithdrawalStatus) (*model.WithdrawalRequest, error) {
	rows, err := client.dbpool.Query(context.Background(),
		"SELECT account_id, wallet_public_key, timestamp, status FROM withdrawal_requests WHERE account_id = $1 AND status = $2", accountId, int(status))

	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var request model.WithdrawalRequest
			err = rows.Scan(&request.AccountId, &request.WalletPublicKey, &request.Timestamp, &request.Status)
			return &request, nil
		}
	}
	return nil, err
}

func (client *PostgresDb) GetWithdrawalRequests(status model.WithdrawalStatus) ([]model.WithdrawalRequest, error) {
	var result []model.WithdrawalRequest

	rows, err := client.dbpool.Query(context.Background(),
		"SELECT account_id, wallet_public_key, timestamp, status, lamports FROM withdrawal_requests WHERE status = $1", int(status))

	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var request model.WithdrawalRequest
			err = rows.Scan(&request.AccountId, &request.WalletPublicKey, &request.Timestamp, &request.Status, &request.Lamports)
			result = append(result, request)
		}
	}
	return result, err
}
