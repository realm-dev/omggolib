package model

type Wallet struct {
	PublicKey string
	AccountId int64
	SecretKey string
	IsPrimary bool
	Timestamp int64
}

type LuckyWallet struct {
	Wallet Wallet
	Hash   string
}
