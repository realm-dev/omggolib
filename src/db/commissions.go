package db

import (
	"context"
	"fmt"
	"github.com/realm-dev/omggolib/src/model"
	"os"
)

func (client *PostgresDb) InsertCommission(commission model.Commission) error {
	commandTag, err := client.dbpool.Exec(context.Background(), "INSERT INTO commissions ("+
		"hash, account_id, ref_account_commission, paid_lamports, timestamp, status, account_public_key, account_fee_public_key,"+
		"token_public_key, mcap, operation_type, volume_sol) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		commission.Hash, commission.AccountId, commission.RefAccountCommission, commission.PaidLamports, commission.Timestamp,
		commission.Status, commission.AccountPublicKey, commission.AccountFeePublicKey, commission.TokenPublicKey,
		commission.MCap, commission.OperationType, commission.VolumeSol)
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

func (client *PostgresDb) GetCommission(hash string) (*model.Commission, error) {
	query := fmt.Sprintf("SELECT "+
		"hash, account_id, ref_account_commission, paid_lamports, timestamp, status, account_public_key, account_fee_public_key, "+
		"token_public_key, mcap, operation_type, volume_sol "+
		"FROM commissions WHERE hash = '%s'", hash)

	rows, err := client.dbpool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var commission model.Commission
		if err := rows.Scan(&commission.Hash, &commission.AccountId, &commission.RefAccountCommission, &commission.PaidLamports,
			&commission.Timestamp, &commission.Status, &commission.AccountPublicKey, &commission.AccountFeePublicKey,
			&commission.TokenPublicKey, &commission.MCap, &commission.OperationType, &commission.VolumeSol); err != nil {
			return nil, err
		}
		return &commission, nil
	}
	return nil, nil
}

func (client *PostgresDb) GetCommissions(accountId int64, status int32) ([]model.Commission, error) {
	var commissions []model.Commission

	rows, err := client.dbpool.Query(context.Background(),
		"SELECT "+
			"hash, account_id, ref_account_commission, paid_lamports, timestamp, status, account_public_key, account_fee_public_key, "+
			"token_public_key, mcap, operation_type, volume_sol "+
			"FROM commissions WHERE account_id = $1 AND status = $2", accountId, status)
	if err != nil {
		return commissions, err
	}

	defer rows.Close()

	for rows.Next() {
		var commission model.Commission
		if err := rows.Scan(&commission.Hash, &commission.AccountId, &commission.RefAccountCommission, &commission.PaidLamports,
			&commission.Timestamp, &commission.Status, &commission.AccountPublicKey, &commission.AccountFeePublicKey,
			&commission.TokenPublicKey, &commission.MCap, &commission.OperationType, &commission.VolumeSol); err != nil {
			return commissions, err
		}
		commissions = append(commissions, commission)
	}
	return commissions, nil
}

func (client *PostgresDb) GetPaybackCommissions(accountId int64, status model.CommissionStatus) ([]model.PaybackCommission, error) {
	var commissions []model.PaybackCommission

	rows, err := client.dbpool.Query(context.Background(),
		"SELECT hash, round(paid_lamports * ref_account_commission / 100.0) as commission "+
			"FROM commissions WHERE account_id = $1 AND status = $2", accountId, int32(status))
	if err != nil {
		return commissions, err
	}

	defer rows.Close()

	for rows.Next() {
		var commission model.PaybackCommission
		if err := rows.Scan(&commission.Hash, &commission.Lamportds); err != nil {
			return commissions, err
		}
		commissions = append(commissions, commission)
	}
	return commissions, nil
}

func (client *PostgresDb) GetTotalCommission(accountId int64, status model.CommissionStatus) (int64, error) {
	var sum int64 = 0
	query := fmt.Sprintf("SELECT round(sum(paid_lamports * ref_account_commission / 100.0)) as commission "+
		"FROM commissions WHERE account_id = %d AND status = %d", accountId, status)

	err := client.dbpool.QueryRow(context.Background(), query).Scan(&sum)
	if err != nil {
		fmt.Errorf("no commission found, error: %v, query: %s", err)
	}
	return sum, nil
}

func (client *PostgresDb) GetLastMCap(accountPublicKey string, tokenPublicKey string) (int64, error) {
	var mcap int64 = 0

	query := fmt.Sprintf("SELECT "+
		"a.mcap FROM commissions AS a "+
		"INNER JOIN "+
		"(SELECT MAX(timestamp) timestamp FROM commissions "+
		"WHERE token_public_key = '%s' AND account_public_key = '%s' AND operation_type = 0) b "+
		"ON a.timestamp = b.timestamp", tokenPublicKey, accountPublicKey)

	err := client.dbpool.QueryRow(context.Background(), query).Scan(&mcap)
	return mcap, err
}

func (client *PostgresDb) UpdateCommissionStatus(hash string, status model.CommissionStatus) error {
	query := fmt.Sprintf("UPDATE commissions SET status = %d WHERE hash = '%s'", int(status), hash)
	_, err := client.dbpool.Exec(context.Background(), query)

	return err
}
