// Package server
package server

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sftx/bank-api/core/ports"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	accountsHandler     ports.IAccountsHandler
	transactionsHandler ports.ITransactionsHandler
}

func NewServer(
	accountsHandler ports.IAccountsHandler,
	transactionsHandler ports.ITransactionsHandler) *Server {
	return &Server{
		accountsHandler:     accountsHandler,
		transactionsHandler: transactionsHandler,
	}
}

func (s *Server) Serve() error {
	version := fmt.Sprintf("/api/%s/", os.Getenv("API_VERSION"))
	port := os.Getenv("API_PORT")
	// SETTING LOGGER
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true})
	router := fiber.New()
	router.Use(recover.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "GET,POST,OPTIONS,HEAD,PUT,DELETE,PATCH",
	}))

	api := router.Group(version)

	////////////// HEALTH SECTION ///////////////////

	api.Get("health", s.accountsHandler.HealthCheck)
	api.Post("accounts", s.accountsHandler.CreateAccount)
	api.Get("accounts", s.accountsHandler.GetAccount)
	api.Post("deposit", s.transactionsHandler.HandleDepositRequest)
	api.Post("withdraw", s.transactionsHandler.HandleWithdrawRequest)
	return router.Listen(fmt.Sprintf(":%s", port))
}
