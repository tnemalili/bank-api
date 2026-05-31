package models

import (
	"fmt"
	"time"
)

type MessageEvent struct {
	Message string `json:"message"`
	Topic   string `json:"topicArn"`
}

type WithdrawalEvent struct {
	Amount    float64 `json:"amount"`
	AccountID int64   `json:"accountId"`
	Status    string  `json:"status"`
}

type WithdrawalResult struct {
	NewBalance Amount    `json:"newBalance" gorm:"embedded"`
	EventID    string    `json:"eventId"`
	Replayed   bool      `json:"replayed"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
	Success    bool      `json:"success"` // true if this was a duplicate idempotency key
}

type DepositResult struct {
	NewBalance Amount    `json:"newBalance" gorm:"embedded"`
	EventID    string    `json:"eventId"`
	Replayed   bool      `json:"replayed"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
	Success    bool      `json:"success"` // true if this was a duplicate idempotency key
}

func (e WithdrawalEvent) toJSON() string {
	// BUG M5: fragile manual JSON, no escaping, amount quoted as a string.
	return fmt.Sprintf(`{"amount":"%f","accountId":%d,"status":"%s"}`,
		e.Amount, e.AccountID, e.Status)
}

func NewWithdrawalEvent(amount float64, accountID int64, status string) WithdrawalEvent {
	return WithdrawalEvent{
		Amount:    amount,
		AccountID: accountID,
		Status:    status,
	}
}

func NewMessageEvent(message, topic string) MessageEvent {
	return MessageEvent{
		Message: message,
		Topic:   topic,
	}
}
