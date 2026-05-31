package services

import (
	"github.com/sftx/bank-api/core/models"
	"github.com/sftx/bank-api/core/ports"
)

type AccountsService struct {
	repo ports.IAccountsRepository
}

// SetBalance implements [ports.IAccountsService].
func (a *AccountsService) SetBalance(accountID string, newBalance float64) error {
	
	return a.repo.SetBalance(accountID, newBalance)
}

// GetAccount implements [ports.IAccountsService].
func (a *AccountsService) GetAccount(accountID string) (models.Account, error) {
	
	return a.repo.GetAccount(accountID)
}

// CreateAccount implements [ports.IAccountsService].
func (a *AccountsService) CreateAccount(acc models.Account) (models.Account, error) {
	
	return a.repo.CreateAccount(acc)
}

// Get feedback in case we are not implementing the interface correctly
var _ ports.IAccountsService = (*AccountsService)(nil)

func NewAccountsService(repo ports.IAccountsRepository) *AccountsService {

	return &AccountsService{repo: repo}
}
