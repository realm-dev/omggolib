package model

type CommissionStatus int32

const (
	CS_Paid CommissionStatus = iota
	CS_Paidback
)

type OperationType int32

const (
	OT_Buy OperationType = iota
	OT_Sell
)

type PaybackCommission struct {
	Hash      string
	Lamportds int64
}

type Commission struct {
	Hash                 string
	AccountId            int64
	RefAccountCommission int32
	PaidLamports         int64
	Timestamp            int64
	Status               CommissionStatus
	AccountPublicKey     string
	AccountFeePublicKey  string
	TokenPublicKey       string
	MCap                 int64
	OperationType        OperationType
	VolumeSol            float64
}
