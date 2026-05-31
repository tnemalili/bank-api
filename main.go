package main

import (
	"github.com/sftx/bank-api/core/services"
	"github.com/sftx/bank-api/handlers"
	"github.com/sftx/bank-api/repository/relationaldb"
	"github.com/sftx/bank-api/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	// SETTING LOGGER
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true})
	// DB Client
	dbClient := relationaldb.NewDBClient()
	
	// Repositories
	accountRepo := relationaldb.NewAccountManager(dbClient)

	transactionRepo := relationaldb.NewTransactionManager(accountRepo, dbClient)
	// Services
	accountServices := services.NewAccountsService(accountRepo)
	transactionServices := services.NewTransactionsService(transactionRepo)
	// Handlers
	accountsHandler := handlers.NewAccountsHandler(accountServices)
	transactionsHandler := handlers.NewTransactionsHandler(transactionServices)
	// Server
	newServer := server.NewServer(accountsHandler, transactionsHandler)
	err := newServer.Serve()
	log.Infof(err.Error())
}
