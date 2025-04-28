package models

import "github.com/google/uuid"

type Account struct {
	AccountID int `json:"account_id"`
	UserID    int `json:"user_id"`
	Balance   int `json:"balance"` // stored in rupees
}

type CreateAccountRequest struct {
	Balance float64 `json:"balance"`
}

type FundTransactionRequest struct {
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	TransactionType string  `json:"transactionType" binding:"required,oneof=deposit withdraw"`
}

// type FundTransactionResponse struct {
// }

type TransactionRequestKafka struct {
	UserId          int       `json:"userId"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transactionType"`
	RequestId       uuid.UUID `json:"requestId"`
	TransactionTime int64     `json:"transactionTime"`
}
