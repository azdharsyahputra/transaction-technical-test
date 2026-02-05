package router

import (
	"github.com/gin-gonic/gin"

	"transaction-technical-test/internal/handler"
)

func RegisterRoutes(
	r *gin.Engine,
	txHandler *handler.TransactionHandler,
	dashboardHandler *handler.DashboardHandler,
) {
	api := r.Group("/api")

	// Transaction routes
	transactions := api.Group("/transactions")
	{
		transactions.POST("", txHandler.Create)
		transactions.GET("", txHandler.GetAll)
		transactions.GET("/:id", txHandler.GetByID)
		transactions.PUT("/:id", txHandler.UpdateStatus)
		transactions.DELETE("/:id", txHandler.Delete)
	}

	// Dashboard routes
	dashboard := api.Group("/dashboard")
	{
		dashboard.GET("/summary", dashboardHandler.Summary)
	}
}
