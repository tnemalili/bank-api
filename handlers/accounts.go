// Package handler provides the public facing interface.
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sftx/bank-api/core/ports"
)

type AccountsHandler struct {
	service ports.IAccountsService
}

// CreateAccount implements [ports.IAccountsHandler].
func (a *AccountsHandler) CreateAccount(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// GetAccount implements [ports.IAccountsHandler].
func (a *AccountsHandler) GetAccount(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// HealthCheck implements [ports.IAccountsHandler].
func (a *AccountsHandler) HealthCheck(ctx *fiber.Ctx) error {
	
	return ctx.JSON(fiber.Map{"status": "ok"})
}


func NewAccountsHandler(service ports.IAccountsService) *AccountsHandler {

	return &AccountsHandler{service: service}
}

var _ ports.IAccountsHandler = (*AccountsHandler)(nil)
