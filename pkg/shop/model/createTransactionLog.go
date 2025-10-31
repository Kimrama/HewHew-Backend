package model

type TransactionLog struct {
	TransactionLogID string  `json:"transactionlog_id"`
	TargetUserID     string  `json:"targetuser_id"`
	OrderID          string  `json:"order_id"`
	Detail           string  `json:"detail"`
	Amount           float64 `json:"amount"`
	TimeStamp        string  `json:"timestamp"`
}
