package relationaldb

import (
	"github.com/sftx/bank-api/core/models"
	"github.com/sftx/bank-api/core/ports"
	"gorm.io/gorm"
)

type accountManager struct {
	dbClient *gorm.DB
}

// SetBalance implements [ports.IAccountsRepository].
func (am *accountManager) SetBalance(accountID string, newBalance float64) error {
	var account models.Account
	if err := am.dbClient.First(&account, "id = ?", accountID).Error; err != nil {
		return err
	}
	account.SetBalance(newBalance)
	if err := am.dbClient.Save(&account).Error; err != nil {
		return err
	}
	return nil
}

// CreateAccount implements [ports.IAccountsRepository].
func (am *accountManager) CreateAccount(req models.Account) (models.Account, error) {
	err := am.dbClient.Create(&req).Error

	if err != nil {
		return models.Account{}, err
	}
	return req, nil
}

// GetAccount implements [ports.IAccountsRepository].
func (am *accountManager) GetAccount(accountID string) (models.Account, error) {
	var account models.Account
	if err := am.dbClient.First(&account, "id = ?", accountID).Error; err != nil {
		return models.Account{}, err
	}
	
	return account, nil
}

func NewAccountManager(dbClient *gorm.DB) *accountManager {
	dbClient.AutoMigrate(&models.Account{})
	return &accountManager{dbClient: dbClient}
}

var _ ports.IAccountsRepository = (*accountManager)(nil)
