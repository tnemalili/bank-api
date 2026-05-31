// Package handler provides the public facing interface.
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sftx/bank-api/core/models"
	"github.com/sftx/bank-api/core/ports"
)

type AccountsHandler struct {
	service ports.IAccountsService
}

// CreateAccount implements [ports.IAccountsHandler].
func (a *AccountsHandler) CreateAccount(ctx *fiber.Ctx) error {
	
	var req models.CreateAccountRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		errMsg := "invalid request body"
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMsg})
	}
	
	newAccount := models.NewAccount(req)
	acc, err := a.service.CreateAccount(newAccount)
	if err != nil {
		errMsg := "failed to create account"
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
	}
	return ctx.Status(fiber.StatusCreated).JSON(acc)
}

// GetAccount implements [ports.IAccountsHandler].
func (a *AccountsHandler) GetAccount(ctx *fiber.Ctx) error {
	accountID := ctx.Request().Header.Peek("X-Account-Id")
	if accountID == nil {
		errMsg := "missing X-Account-Id header"
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMsg})
	}
	
	acc, err := a.service.GetAccount(string(accountID))
	if err != nil {
		errMsg := "failed to retrieve account"
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
	}
	return ctx.JSON(acc)
}

// HealthCheck implements [ports.IAccountsHandler].
func (a *AccountsHandler) HealthCheck(ctx *fiber.Ctx) error {

	return ctx.JSON(fiber.Map{"status": "ok"})
}


func NewAccountsHandler(service ports.IAccountsService) *AccountsHandler {

	return &AccountsHandler{service: service}
}

var _ ports.IAccountsHandler = (*AccountsHandler)(nil)
