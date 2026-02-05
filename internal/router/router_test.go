package router

import (
	"testing"

	"transaction-technical-test/internal/handler"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	r := gin.New()

	RegisterRoutes(
		r,
		&handler.TransactionHandler{},
		&handler.DashboardHandler{},
	)
}
