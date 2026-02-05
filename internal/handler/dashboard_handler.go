package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"transaction-technical-test/internal/service"
)

type DashboardHandler struct {
	service *service.DashboardService
	logger  *zap.Logger
}

func NewDashboardHandler(s *service.DashboardService, logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		service: s,
		logger:  logger,
	}
}

func (h *DashboardHandler) Summary(c *gin.Context) {
	summary, err := h.service.GetSummary()
	if err != nil {
		h.logger.Error("failed to get dashboard summary", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("dashboard summary retrieved")
	c.JSON(http.StatusOK, summary)
}
