package models

import (
	"fmt"
	"sync/atomic"
	"time"
)

var accountIDCounter uint64 = uint64(time.Now().UnixNano() % 10000000000)

type CreateAccountRequest struct {
	AccountHolder    string  `json:"accountHolder"`
	InitiationAmount float64 `json:"initiationAmount"`
	Currency         string  `json:"currency"`
}

type WithdrawRequest struct {
	AccountID      string  `json:"accountId"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotencyKey"`
}

type DepositRequest struct {
	AccountID      string  `json:"accountId"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotencyKey"`
}

type Amount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type Account struct {
	ID            string  `json:"id" gorm:"primaryKey"`
	AccountHolder string  `json:"accountHolder"`
	AccountNumber string  `json:"accountNumber"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
}

// Getters for WithdrawRequest

func (wr *WithdrawRequest) GetAccountID() string {
	return wr.AccountID
}

func (wr *WithdrawRequest) GetAmount() float64 {
	return wr.Amount
}

func (wr *WithdrawRequest) GetCurrency() string {
	return wr.Currency
}

func (wr *WithdrawRequest) GetIdempotencyKey() string {
	return wr.IdempotencyKey
}

// Getters for DepositRequest

func (dr *DepositRequest) GetAccountID() string {
	return dr.AccountID
}

func (dr *DepositRequest) GetAmount() float64 {
	return dr.Amount
}

func (dr *DepositRequest) GetCurrency() string {
	return dr.Currency
}

func (dr *DepositRequest) GetIdempotencyKey() string {
	return dr.IdempotencyKey
}

// Getters for Account

func (a *Account) GetAccountID() string {
	return a.ID
}

func (a *Account) GetBalance() float64 {
	return a.Balance
}

func (a *Account) GetCurrency() string {
	return a.Currency
}

func (a *Account) GetAccountHolder() string {
	return a.AccountHolder
}

// Setters for Account

func (a *Account) SetBalance(newBalance float64) {
	a.Balance = newBalance
}

func NewAccount(req CreateAccountRequest) Account {

	return Account{
		ID:            generateAccountID(),
		AccountHolder: req.AccountHolder,
		Balance:       req.InitiationAmount,
		Currency:      req.Currency,
		AccountNumber: generateAccountID(),
	}
}

func generateAccountID() string {
	accountID := atomic.AddUint64(&accountIDCounter, 1) % 10000000000
	return fmt.Sprintf("%010d", accountID)
}
