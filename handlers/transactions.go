package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sftx/bank-api/core/models"
	"github.com/sftx/bank-api/core/ports"
)

type TransactionsHandler struct {
	service ports.ITransactionsService
}

// HandleDepositRequest implements [ports.ITransactionsHandler].
func (t *TransactionsHandler) HandleDepositRequest(ctx *fiber.Ctx) error {
	
	var req models.DepositRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		errMsg := "invalid request body"
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMsg})
	}
	
	result := t.service.Deposit(req)
	return ctx.JSON(result)
}

// HandleWithdrawRequest implements [ports.ITransactionsHandler].
func (t *TransactionsHandler) HandleWithdrawRequest(ctx *fiber.Ctx) error {
	
	var req models.WithdrawRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		errMsg := "invalid request body"
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMsg})
	}
	
	result := t.service.Withdraw(req)
	return ctx.JSON(result)	
}

func NewTransactionsHandler(service ports.ITransactionsService) *TransactionsHandler {
	return &TransactionsHandler{service: service}
}

var _ ports.ITransactionsHandler = (*TransactionsHandler)(nil)
