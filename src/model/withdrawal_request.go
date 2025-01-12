package model

type WithdrawalStatus int32

const (
	WS_Requested WithdrawalStatus = iota
	WS_PaidOut
)

type WithdrawalRequest struct {
	AccountId          int64
	WalletPublicKey    string
	Timestamp          int64
	Status             WithdrawalStatus
	CalculatedLamports int64
	PaidLamports       int64
}
