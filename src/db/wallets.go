package db

import (
	"context"
	"fmt"
	"github.com/realm-dev/omggolib/src/model"
	"os"
)

func (client *PostgresDb) InsertWallet(wallet model.Wallet) error {
	commandTag, err := client.dbpool.Exec(context.Background(), "INSERT INTO wallets (public_key, account_id, secret_key, is_primary) VALUES ($1, $2, $3, $4)",
		wallet.PublicKey, wallet.AccountId, wallet.SecretKey, wallet.IsPrimary)
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

func (client *PostgresDb) GetWalletsCount(accountId int64) (int, error) {
	accountWalletsCount := 0
	err := client.dbpool.QueryRow(context.Background(),
		"SELECT count(*) FROM wallets WHERE account_id = $1",
		accountId).Scan(&accountWalletsCount)
	return accountWalletsCount, err
}

func (client *PostgresDb) SetPrimaryKey(accountId int64, newPrimaryPublicKey string) error {
	_, err := client.dbpool.Exec(context.Background(),
		"UPDATE wallets SET is_primary = false WHERE account_id = $1",
		accountId)

	if err == nil {
		_, err = client.dbpool.Exec(context.Background(),
			"UPDATE wallets SET is_primary = TRUE WHERE account_id = $1 AND public_key = $2",
			accountId, newPrimaryPublicKey)
	}
	return err
}

func (client *PostgresDb) GetPrimaryWallet(accountId int64) (*model.Wallet, error) {
	rows, err := client.dbpool.Query(context.Background(),
		"SELECT public_key, account_id, secret_key, is_primary FROM wallets WHERE account_id = $1 AND is_primary = TRUE", accountId)

	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var primaryWallet model.Wallet
			err = rows.Scan(&primaryWallet.PublicKey, &primaryWallet.AccountId, &primaryWallet.SecretKey, &primaryWallet.IsPrimary)
			return &primaryWallet, nil
		}
	}
	return nil, err
}

func (client *PostgresDb) GetWallets(accountId int64) ([]model.Wallet, error) {
	var wallets []model.Wallet

	rows, err := client.dbpool.Query(context.Background(),
		"SELECT public_key, account_id, secret_key, is_primary FROM wallets WHERE account_id = $1 ORDER BY public_key ASC", accountId)

	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var wallet model.Wallet
			if err = rows.Scan(&wallet.PublicKey, &wallet.AccountId, &wallet.SecretKey, &wallet.IsPrimary); err != nil {
				return wallets, err
			}
			wallets = append(wallets, wallet)
		}
	}

	return wallets, err
}
