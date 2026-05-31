package relationaldb

import (
	"errors"
	"strings"

	"github.com/sftx/bank-api/core/models"
	"github.com/sftx/bank-api/core/ports"
	"gorm.io/gorm"
)

type accountManager struct {
	dbClient *gorm.DB
}

// SetBalance implements [ports.IAccountsRepository].
func (am *accountManager) SetBalance(accountID string, newBalance float64) error {
	account, err := am.findAccountByIdentifier(accountID)
	if err != nil {
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
	account, err := am.findAccountByIdentifier(accountID)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (am *accountManager) findAccountByIdentifier(accountID string) (models.Account, error) {
	var account models.Account

	err := am.dbClient.Where("id = ? OR account_number = ?", accountID, accountID).First(&account).Error
	if err == nil {
		return account, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Account{}, err
	}

	legacyErr := am.dbClient.Where("account_id = ?", accountID).First(&account).Error
	if legacyErr == nil {
		return account, nil
	}

	if strings.Contains(legacyErr.Error(), "no such column") {
		return models.Account{}, err
	}

	return models.Account{}, legacyErr
}

func NewAccountManager(dbClient *gorm.DB) *accountManager {
	dbClient.AutoMigrate(&models.Account{})
	return &accountManager{dbClient: dbClient}
}

var _ ports.IAccountsRepository = (*accountManager)(nil)
