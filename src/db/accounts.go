package db

import (
	"context"
	"fmt"
	"github.com/realm-dev/omggolib/src/model"
	"os"
)

func (client *PostgresDb) InsertAccount(account model.Account) error {
	commandTag, err := client.dbpool.Exec(context.Background(), "INSERT INTO accounts ("+
		"account_id, alias_id, account_role, ref_account_id, affiliate_level, commission_discount, buy_token_volume, sell_percent, slippage, priority_fee, username) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		account.AccountId, account.AliasId, account.AccountRole, account.RefAccountId, account.AffiliateLevel, account.CommissionDiscount,
		account.BuyTokenVolume, account.SellPercent, account.Slippage, account.PriorityFee, account.Username)
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

func (client *PostgresDb) GetAccount(accountId int64) (*model.Account, error) {
	var account model.Account
	err := client.dbpool.QueryRow(context.Background(),
		"SELECT "+
			"account_id, alias_id, account_role, ref_account_id, affiliate_level, commission_discount, buy_token_volume, sell_percent, slippage, priority_fee, username "+
			"FROM accounts WHERE account_id = $1",
		accountId).Scan(&account.AccountId, &account.AliasId, &account.AccountRole, &account.RefAccountId, &account.AffiliateLevel, &account.CommissionDiscount,
		&account.BuyTokenVolume, &account.SellPercent, &account.Slippage, &account.PriorityFee, &account.Username)
	return &account, err
}

func (client *PostgresDb) GetAccountsByRef(refAccountId int64) ([]model.Account, error) {
	var accounts []model.Account

	rows, err := client.dbpool.Query(context.Background(),
		"SELECT "+
			"account_id, alias_id, account_role, ref_account_id, affiliate_level, commission_discount, buy_token_volume, sell_percent, slippage, priority_fee, username "+
			"FROM Accounts WHERE ref_account_id = $1", refAccountId)
	if err != nil {
		return accounts, err
	}

	defer rows.Close()

	for rows.Next() {
		var account model.Account
		if err := rows.Scan(&account.AccountId, &account.AliasId, &account.AccountRole, &account.RefAccountId, &account.AffiliateLevel, &account.CommissionDiscount,
			&account.BuyTokenVolume, &account.SellPercent, &account.Slippage, &account.PriorityFee, &account.Username); err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (client *PostgresDb) GetAccountIdsByRef(refAccountId int64) ([]int64, error) {
	var accounts []int64

	rows, err := client.dbpool.Query(context.Background(),
		"SELECT account_id FROM Accounts WHERE ref_account_id = $1", refAccountId)
	if err != nil {
		return accounts, err
	}

	defer rows.Close()

	for rows.Next() {
		var accountId int64
		if err := rows.Scan(&accountId); err != nil {
			return accounts, err
		}
		accounts = append(accounts, accountId)
	}
	return accounts, nil
}

func (client *PostgresDb) GetAccountsCountByRef(refAccountId int64) (int64, error) {
	var count int64 = 0
	err := client.dbpool.QueryRow(context.Background(),
		"SELECT count(*) "+
			"FROM accounts WHERE ref_account_id = $1",
		refAccountId).Scan(&count)
	return count, err
}

func (client *PostgresDb) GetAccountByAlias(aliasId string) (*model.Account, error) {
	var account model.Account
	err := client.dbpool.QueryRow(context.Background(),
		"SELECT "+
			"account_id, alias_id, account_role, ref_account_id, affiliate_level, commission_discount, buy_token_volume, sell_percent, slippage, priority_fee, username "+
			"FROM accounts WHERE alias_id = $1", aliasId).
		Scan(&account.AccountId, &account.AliasId, &account.AccountRole, &account.RefAccountId, &account.AffiliateLevel, &account.CommissionDiscount,
			&account.BuyTokenVolume, &account.SellPercent, &account.Slippage, &account.PriorityFee, &account.Username)
	return &account, err
}
