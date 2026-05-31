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
func (t *transactionManager) Deposit(req models.DepositRequest) models.DepositResult {
	account, err := t.accountRepo.GetAccount(req.AccountID)

	if err != nil {
		return models.DepositResult{
			Status:  "failed",
			Message: "Account not found",
		}
	}

	account.SetBalance(account.GetBalance() + req.GetAmount())

	err = t.dbClient.Save(&account).Error

	if err != nil {
		return models.DepositResult{
			Status:  "failed",
			Success: false,
			Message: "Failed to update account balance",
		}
	}

	return models.DepositResult{
		Status:  "success",
		Success: true,
		Message: "Deposit successful",
	}
}

// Withdraw implements [ports.ITransactionsRepository].
func (t *transactionManager) Withdraw(req models.WithdrawRequest) models.WithdrawalResult {

	account, err := t.accountRepo.GetAccount(req.AccountID)

	if err != nil {
		return models.WithdrawalResult{
			Status:  "failed",
			Message: "Account not found",
		}
	}

	if account.GetBalance() < req.GetAmount() {
		return models.WithdrawalResult{
			Status:  "failed",
			Success: false,
			Message: "Insufficient funds",
		}
	}

	account.SetBalance(account.GetBalance() - req.GetAmount())

	err = t.dbClient.Save(&account).Error

	if err != nil {
		return models.WithdrawalResult{
			Status:  "failed",
			Success: false,
			Message: "Failed to update account balance",
		}
	}

	return models.WithdrawalResult{
		Status:  "success",
		Success: true,
		Message: "Withdrawal successful",
	}
}

func NewTransactionManager(accountRepo ports.IAccountsRepository, dbClient *gorm.DB) *transactionManager {
	dbClient.AutoMigrate(&models.WithdrawalResult{})
	return &transactionManager{accountRepo: accountRepo, dbClient: dbClient}
}

var _ ports.ITransactionsRepository = (*transactionManager)(nil)
