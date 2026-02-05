package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"transaction-technical-test/internal/config"
	"transaction-technical-test/internal/handler"
	"transaction-technical-test/internal/repository"
	"transaction-technical-test/internal/router"
	"transaction-technical-test/internal/service"
)

func main() {
	// Init database
	db := config.InitDB()

	// Init logger
	logger := config.InitLogger()
	defer logger.Sync()

	// Repository
	transactionRepo := repository.NewTransactionRepository(db)

	// Service
	transactionService := service.NewTransactionService(transactionRepo)
	dashboardService := service.NewDashboardService(transactionRepo)

	// Handler
	transactionHandler := handler.NewTransactionHandler(transactionService, logger)
	dashboardHandler := handler.NewDashboardHandler(dashboardService, logger)

	// Router
	r := gin.Default()
	router.RegisterRoutes(r, transactionHandler, dashboardHandler)

	log.Println("server running on :8080")
	log.Fatal(r.Run(":8080"))
}
