// Package services
package services

import (
	"github.com/sftx/bank-api/core/models"
	"github.com/sftx/bank-api/core/ports"
)

type TransactionsService struct {
	repo ports.ITransactionsRepository
}

// Deposit implements [ports.ITransactionsService].
func (t *TransactionsService) Deposit(req models.DepositRequest) models.DepositResult {
	
	return t.repo.Deposit(req)
}

// Withdraw implements [ports.ITransactionsService].
func (t *TransactionsService) Withdraw(req models.WithdrawRequest) models.WithdrawalResult {

	return t.repo.Withdraw(req)
}

// Get feedback in case we are not implementing the interface correctly
var _ ports.ITransactionsService = (*TransactionsService)(nil)

func NewTransactionsService(repo ports.ITransactionsRepository) *TransactionsService {

	return &TransactionsService{repo: repo}
}
