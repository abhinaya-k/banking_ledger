package models

import "github.com/google/uuid"

type Account struct {
	AccountID int   `json:"account_id"`
	UserID    int   `json:"user_id"`
	Balance   int64 `json:"balance"` // stored in rupees
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

type TransactionCollection struct {
	UserId            int       `bson:"userId"`
	Amount            float64   `bson:"amount"`
	TransactionType   string    `bson:"transactionType"`
	TransactionStatus string    `bson:"transactionStatus"`
	TransactionMsg    string    `bson:"transactionMessage"`
	RequestId         uuid.UUID `bson:"requestId"`
	TransactionTime   int64     `bson:"transactionTime"`
}

type GetTransactionHistoryRequest struct {
	Filters *struct {
		TransactionType *string `json:"transactionType,omitempty" binding:"oneof=deposit withdraw"`
		StartTime       *int64  `json:"startTime,omitempty"`
		EndTime         *int64  `json:"endTime,omitempty"`
	} `json:"filters,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type TransactionHistory struct {
	UserId            int     `json:"userId"`
	FirstName         string  `json:"fistName"`
	LastName          string  `json:"lastName"`
	Amount            float64 `json:"amount"`
	TransactionType   string  `json:"transactionType"`
	TransactionTime   int64   `json:"transactionTime"`
	TransactionStatus string  `bson:"transactionStatus"`
	TransactionMsg    string  `bson:"transactionMessage"`
}

type GetTransactionHistoryResponse struct {
	TransactionHistory []TransactionHistory `json:"transactionHistory"`
	Pagination         Pagination           `json:"pagination"`
}

type Pagination struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}
