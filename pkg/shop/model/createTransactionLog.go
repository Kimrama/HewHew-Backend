package model

type TransactionLog struct {
	TransactionLogID string  `json:"transactionlogid"`
	TargetUserID     string  `json:"targetuserid"`
	OrderID          string  `json:"orderid"`
	Detail           string  `json:"detail"`
	Amount           float64 `json:"amount"`
	TimeStamp        string  `json:"timestamp"`
}
