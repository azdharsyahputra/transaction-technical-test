package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"transaction-technical-test/internal/domain"
	"transaction-technical-test/internal/service"
)

type TransactionHandler struct {
	service *service.TransactionService
	logger  *zap.Logger
}

func NewTransactionHandler(s *service.TransactionService, logger *zap.Logger) *TransactionHandler {
	return &TransactionHandler{
		service: s,
		logger:  logger,
	}
}

type CreateTransactionRequest struct {
	UserID uint    `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type UpdateStatusRequest struct {
	Status domain.TransactionStatus `json:"status" binding:"required"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid create transaction request",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.service.Create(req.UserID, req.Amount)
	if err != nil {
		h.logger.Error("failed to create transaction",
			zap.Uint("user_id", req.UserID),
			zap.Float64("amount", req.Amount),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("transaction created",
		zap.Uint("transaction_id", tx.ID),
		zap.Uint("user_id", tx.UserID),
		zap.Float64("amount", tx.Amount),
	)

	c.JSON(http.StatusCreated, tx)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("invalid transaction id",
			zap.String("id", idStr),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tx, err := h.service.GetByID(uint(id))
	if err != nil {
		if err == domain.ErrTransactionNotFound {
			h.logger.Info("transaction not found",
				zap.Uint("transaction_id", uint(id)),
			)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error("failed to get transaction",
			zap.Uint("transaction_id", uint(id)),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("transaction retrieved",
		zap.Uint("transaction_id", tx.ID),
	)

	c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	var filter domain.TransactionFilter

	if userID := c.Query("user_id"); userID != "" {
		id, err := strconv.Atoi(userID)
		if err != nil {
			h.logger.Warn("invalid user_id query",
				zap.String("user_id", userID),
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		uid := uint(id)
		filter.UserID = &uid
	}

	if status := c.Query("status"); status != "" {
		s := domain.TransactionStatus(status)
		filter.Status = &s
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter.Limit = limit
	filter.Offset = (page - 1) * limit

	result, err := h.service.GetAll(filter)
	if err != nil {
		h.logger.Error("failed to get transactions",
			zap.Any("filter", filter),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("transactions retrieved",
		zap.Int("count", len(result)),
		zap.Any("filter", filter),
	)

	c.JSON(http.StatusOK, result)
}

func (h *TransactionHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("invalid transaction id", zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid update status request",
			zap.Uint("transaction_id", uint(id)),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		h.logger.Error("failed to update transaction status",
			zap.Uint("transaction_id", uint(id)),
			zap.String("status", string(req.Status)),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("transaction status updated",
		zap.Uint("transaction_id", uint(id)),
		zap.String("status", string(req.Status)),
	)

	c.Status(http.StatusNoContent)
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("invalid transaction id",
			zap.String("id", idStr),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if err == domain.ErrTransactionNotFound {
			h.logger.Warn("transaction not found",
				zap.Uint("transaction_id", uint(id)),
			)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error("failed to delete transaction",
			zap.Uint("transaction_id", uint(id)),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("transaction deleted",
		zap.Uint("transaction_id", uint(id)),
	)

	c.Status(http.StatusNoContent)
}
