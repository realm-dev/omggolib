package model

type CommissionStatus int32

const (
	Paid CommissionStatus = iota
	Paidback
)

type OperationType int32

const (
	Buy OperationType = iota
	Sell
)

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
