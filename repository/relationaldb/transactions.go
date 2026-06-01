package relationaldb

import (
	"github.com/sftx/bank-api/core/models"
	"github.com/sftx/bank-api/core/ports"
	"gorm.io/gorm"
)

type transactionManager struct {
	accountRepo ports.IAccountsRepository
	dbClient    *gorm.DB
}

// Deposit implements [ports.ITransactionsRepository].
func (t *transactionManager) Deposit(req models.DepositRequest) models.TransactionResult {
	
	if req.GetAmount() <= 0 {
		return models.NewTransactionResult(models.TransactionEvent{
			Amount:     req.GetAmount(),
			Currency:   req.GetCurrency(),
			Status:     "failed",
			Message:    "Amount must be greater than zero",
			Success:    false,
			StatusCode: 400,
			NewBalance: 0, // Default to 0 for failed transactions
		})
	}

	account, err := t.accountRepo.GetAccount(req.AccountID)

	if err != nil {
		return models.NewTransactionResult(models.TransactionEvent{
			Amount:     req.GetAmount(),
			Currency:   req.GetCurrency(),
			Status:     "failed",
			Message:    "Account not found",
			Success:    false,
			StatusCode: 404,
			NewBalance: 0, // Default to 0 for failed transactions
		})
	}

	newBalance := account.GetBalance() + req.GetAmount()

	account.SetBalance(newBalance)

	err = t.dbClient.Save(&account).Error

	if err != nil {
		return models.NewTransactionResult(models.TransactionEvent{
			Amount:     req.GetAmount(),
			Currency:   req.GetCurrency(),
			Status:     "failed",
			Message:    "Failed to update account balance",
			Success:    false,
			StatusCode: 500,
			NewBalance: account.GetBalance(),
		})
	}

	return models.NewTransactionResult(models.TransactionEvent{
		Amount:     req.GetAmount(),
		Currency:   req.GetCurrency(),
		NewBalance: newBalance,
		Status:     "success",
		Message:    "Deposit successful",
		Success:    true,
		StatusCode: 200,
	})
}

// Withdraw implements [ports.ITransactionsRepository].
func (t *transactionManager) Withdraw(req models.WithdrawRequest) models.TransactionResult {
	if req.GetAmount() <= 0 {
		return models.NewTransactionResult(models.TransactionEvent{
			Amount:     req.GetAmount(),
			Currency:   req.GetCurrency(),
			Status:     "failed",
			Message:    "Amount must be greater than zero",
			Success:    false,
			StatusCode: 400,
			NewBalance: 0, // Default to 0 for failed transactions
		})
	}

	account, err := t.accountRepo.GetAccount(req.AccountID)

	if err != nil {
		return models.NewTransactionResult(models.TransactionEvent{
			Amount:     req.GetAmount(),
			Currency:   req.GetCurrency(),
			Status:     "failed",
			Message:    "Account not found",
			Success:    false,
			StatusCode: 404,
			NewBalance: 0, // Default to 0 for failed transactions
		})
	}

	if account.GetBalance() < req.GetAmount() {
		return models.NewTransactionResult(models.TransactionEvent{
			Amount:     req.GetAmount(),
			Currency:   req.GetCurrency(),
			NewBalance: account.GetBalance(),
			Status:     "failed",
			Message:    "Insufficient funds",
			Success:    false,
			StatusCode: 400,
		})
	}

	newBalance := account.GetBalance() - req.GetAmount()
	account.SetBalance(newBalance)

	err = t.dbClient.Save(&account).Error

	if err != nil {
		return models.NewTransactionResult(models.TransactionEvent{
			Amount:     req.GetAmount(),
			Currency:   req.GetCurrency(),
			NewBalance: account.GetBalance(),
			Status:     "failed",
			Message:    "Failed to update account balance",
			Success:    false,
			StatusCode: 500,
		})
	}

	return models.NewTransactionResult(models.TransactionEvent{
		Amount:     req.GetAmount(),
		Currency:   req.GetCurrency(),
		NewBalance: newBalance,
		Status:     "success",
		Message:    "Withdrawal successful",
		Success:    true,
		StatusCode: 200,
	})
}

func NewTransactionManager(accountRepo ports.IAccountsRepository, dbClient *gorm.DB) *transactionManager {
	dbClient.AutoMigrate(&models.TransactionResult{})
	return &transactionManager{accountRepo: accountRepo, dbClient: dbClient}
}

var _ ports.ITransactionsRepository = (*transactionManager)(nil)
