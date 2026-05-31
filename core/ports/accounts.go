package ports

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sftx/bank-api/core/models"
)

type IAccountsHandler interface {
	HealthCheck(ctx *fiber.Ctx) error
	CreateAccount(ctx *fiber.Ctx) error
	GetAccount(ctx *fiber.Ctx) error
}

type IAccountsService interface {
	CreateAccount(req models.Account) (models.Account, error)
	GetAccount(accountID string) (models.Account, error)
	SetBalance(accountID string, newBalance float64) error
}

type IAccountsRepository interface {
	CreateAccount(req models.Account) (models.Account, error)
	GetAccount(accountID string) (models.Account, error)
	SetBalance(accountID string, newBalance float64) error
}