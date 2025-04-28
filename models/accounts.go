package models

type Account struct {
	AccountID int `json:"account_id"`
	UserID    int `json:"user_id"`
	Balance   int `json:"balance"` // stored in rupees
}

type CreateAccountRequest struct {
	Balance float64 `json:"balance"`
}

type FundTransactionRequest struct {
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transactionType"`
}

// type FundTransactionResponse struct {
// }
