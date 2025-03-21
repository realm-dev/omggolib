package model

type Role int32
type Tips int32

const (
	System Role = iota
	Admin
	Trader
)

const SystemAccountId = 1

const (
	TipsNone Tips = iota
	Tipsx1
	Tipsx2
	Tipsx3
)

type Account struct {
	AccountId          int64
	AliasId            string
	AccountRole        Role
	RefAccountId       int64
	AffiliateLevel     int32
	CommissionDiscount int32
	BuyTokenVolume     float64
	SellPercent        float64
	Slippage           uint32
	PriorityFee        float64
	Username           string
	ChatId             int64
	JitoTipsBuy        Tips
	JitoTipsSell       Tips
}
