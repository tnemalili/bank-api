package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sftx/bank-api/core/models"
	"github.com/sftx/bank-api/core/ports"
	"github.com/sftx/bank-api/messaging"
	log "github.com/sirupsen/logrus"
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
	// Send transation event
	newMessagingClient := messaging.NewMessagingClient()
	// Publish the transaction event asynchronously to avoid blocking the response
	go func() {
		if err := newMessagingClient.Publish("transaction-topic", result); err != nil {
			log.Errorf("Failed to publish transaction event: %v", err)
		}
	}()
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
	// Send transation event
	newMessagingClient := messaging.NewMessagingClient()
	// Publish the transaction event asynchronously to avoid blocking the response
	go func() {
		if err := newMessagingClient.Publish("transaction-topic", result); err != nil {
			log.Errorf("Failed to publish transaction event: %v", err)
		}
	}()
	return ctx.JSON(result)	
}

func NewTransactionsHandler(service ports.ITransactionsService) *TransactionsHandler {
	return &TransactionsHandler{service: service}
}

var _ ports.ITransactionsHandler = (*TransactionsHandler)(nil)
