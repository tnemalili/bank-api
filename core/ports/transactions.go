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
	Withdraw(req models.WithdrawRequest) models.WithdrawalResult
	Deposit(req models.DepositRequest) models.DepositResult
}

type ITransactionsRepository interface {
	Withdraw(req models.WithdrawRequest) models.WithdrawalResult
	Deposit(req models.DepositRequest) models.DepositResult
}