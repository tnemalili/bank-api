package models

import (
	"fmt"
	"time"
)

type TransactionEvent struct {
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount"`
	NewBalance float64 `json:"newBalance"`
	AccountID  int64   `json:"accountId"`
	Status     string  `json:"status"`
	Replayed   bool    `json:"replayed"`
	Message    string  `json:"message"`
	Success    bool    `json:"success"`
	StatusCode int     `json:"statusCode"` // true if this was a duplicate idempotency key
}

type TransactionResult struct {
	Amount     Amount    `json:"amount" gorm:"embedded"`
	NewBalance Amount    `json:"newBalance" gorm:"embedded"`
	EventID    string    `json:"eventId"`
	Replayed   bool      `json:"replayed"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
	Success    bool      `json:"success"`
	StatusCode int       `json:"statusCode"` // true if this was a duplicate idempotency key
}

func (e TransactionResult) toJSON() string {
	// BUG Major-5: fragile manual JSON, no escaping, amount quoted as a string.
	return fmt.Sprintf(`{"newBalance":{"value":"%f","currency":"%s"},"eventId":"%s","replayed":%t,"status":"%s","message":"%s","createdAt":"%s","success":%t}`,
		e.NewBalance.Value, e.NewBalance.Currency, e.EventID, e.Replayed, e.Status, e.Message, e.CreatedAt.Format(time.RFC3339), e.Success)
}

func (r *TransactionResult) GetStatusCode() int {
	return r.StatusCode
}

func NewTransactionResult(event TransactionEvent) TransactionResult {
	return TransactionResult{
		Amount: Amount{
			Value:    event.Amount,
			Currency: event.Currency,
		},
		NewBalance: Amount{
			Value:    event.NewBalance,
			Currency: event.Currency,
		},
		EventID:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Replayed:  false,
		Status:    event.Status,
		Message:   event.Message,
		CreatedAt: time.Now(),
		Success:   event.Success,
		StatusCode: event.StatusCode,
	}
}