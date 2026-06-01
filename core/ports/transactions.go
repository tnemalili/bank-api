package ports

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sftx/bank-api/core/models"
)

type ITransactionsHandler interface {
	HandleWithdrawRequest(ctx *fiber.Ctx) error
	HandleDepositRequest(ctx *fiber.Ctx) error
}

type ITransactionsService interface {
	Withdraw(req models.WithdrawRequest) models.TransactionResult
	Deposit(req models.DepositRequest) models.TransactionResult
}

type ITransactionsRepository interface {
	Withdraw(req models.WithdrawRequest) models.TransactionResult
	Deposit(req models.DepositRequest) models.TransactionResult
}