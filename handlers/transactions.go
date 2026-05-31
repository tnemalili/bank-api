package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sftx/bank-api/core/ports"
)

type TransactionsHandler struct {
	service ports.ITransactionsService
}

// HandleDepositRequest implements [ports.ITransactionsHandler].
func (t *TransactionsHandler) HandleDepositRequest(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// HandleWithdrawRequest implements [ports.ITransactionsHandler].
func (t *TransactionsHandler) HandleWithdrawRequest(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

func NewTransactionsHandler(service ports.ITransactionsService) *TransactionsHandler {
	return &TransactionsHandler{service: service}
}

var _ ports.ITransactionsHandler = (*TransactionsHandler)(nil)
