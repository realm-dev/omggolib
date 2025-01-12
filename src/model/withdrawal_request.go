package model

type WithdrawalStatus int32

const (
	Requested WithdrawalStatus = iota
)

type WithdrawalRequest struct {
	AccountId       int64
	WalletPublicKey string
	Timestamp       int64
	Status          WithdrawalStatus
	Lamports        int64
}
