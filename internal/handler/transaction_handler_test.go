package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"transaction-technical-test/internal/domain"
	"transaction-technical-test/internal/handler"
	"transaction-technical-test/internal/service"
)

type mockTransactionRepo struct {
	createFn   func(tx *domain.Transaction) error
	findByIDFn func(id uint) (*domain.Transaction, error)
	findAllFn  func(filter domain.TransactionFilter) ([]domain.Transaction, error)
	updateFn   func(tx *domain.Transaction) error
	deleteFn   func(id uint) error
}

func (m *mockTransactionRepo) Create(tx *domain.Transaction) error {
	return m.createFn(tx)
}
func (m *mockTransactionRepo) FindByID(id uint) (*domain.Transaction, error) {
	return m.findByIDFn(id)
}
func (m *mockTransactionRepo) FindAll(filter domain.TransactionFilter) ([]domain.Transaction, error) {
	return m.findAllFn(filter)
}
func (m *mockTransactionRepo) Update(tx *domain.Transaction) error {
	return m.updateFn(tx)
}
func (m *mockTransactionRepo) Delete(id uint) error {
	return m.deleteFn(id)
}

func (m *mockTransactionRepo) TotalSuccessToday() (float64, error) { return 0, nil }
func (m *mockTransactionRepo) AverageAmountPerUser() (float64, error) {
	return 0, nil
}
func (m *mockTransactionRepo) Latest(limit int) ([]domain.Transaction, error) {
	return nil, nil
}
func setupTransactionRouter(repo *mockTransactionRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)

	svc := service.NewTransactionService(repo)

	logger := zap.NewNop()
	h := handler.NewTransactionHandler(svc, logger)

	r := gin.New()
	r.POST("/transactions", h.Create)
	r.GET("/transactions/:id", h.GetByID)
	r.GET("/transactions", h.GetAll)
	r.PUT("/transactions/:id", h.UpdateStatus)
	r.DELETE("/transactions/:id", h.Delete)

	return r
}
func TestTransactionHandler_Create_Success(t *testing.T) {
	repo := &mockTransactionRepo{
		createFn: func(tx *domain.Transaction) error {
			tx.ID = 1
			return nil
		},
	}

	r := setupTransactionRouter(repo)

	body := `{"user_id":1,"amount":1000}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}
func TestTransactionHandler_Create_InvalidBody(t *testing.T) {
	repo := &mockTransactionRepo{}
	r := setupTransactionRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400")
	}
}
func TestTransactionHandler_GetByID_Success(t *testing.T) {
	repo := &mockTransactionRepo{
		findByIDFn: func(id uint) (*domain.Transaction, error) {
			return &domain.Transaction{ID: id}, nil
		},
	}

	r := setupTransactionRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/transactions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
}
func TestTransactionHandler_GetByID_NotFound(t *testing.T) {
	repo := &mockTransactionRepo{
		findByIDFn: func(id uint) (*domain.Transaction, error) {
			return nil, domain.ErrTransactionNotFound
		},
	}

	r := setupTransactionRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/transactions/99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404")
	}
}
func TestTransactionHandler_GetAll_Success(t *testing.T) {
	repo := &mockTransactionRepo{
		findAllFn: func(filter domain.TransactionFilter) ([]domain.Transaction, error) {
			return []domain.Transaction{{ID: 1}, {ID: 2}}, nil
		},
	}

	r := setupTransactionRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
}
func TestTransactionHandler_UpdateStatus_Success(t *testing.T) {
	tx := &domain.Transaction{ID: 1}

	repo := &mockTransactionRepo{
		findByIDFn: func(id uint) (*domain.Transaction, error) {
			return tx, nil
		},
		updateFn: func(tx *domain.Transaction) error {
			return nil
		},
	}

	r := setupTransactionRouter(repo)

	body := map[string]string{"status": string(domain.StatusSuccess)}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/transactions/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204")
	}
}
func TestTransactionHandler_UpdateStatus_Invalid(t *testing.T) {
	repo := &mockTransactionRepo{
		findByIDFn: func(id uint) (*domain.Transaction, error) {
			return &domain.Transaction{ID: id}, nil
		},
	}

	r := setupTransactionRouter(repo)

	req := httptest.NewRequest(
		http.MethodPut,
		"/transactions/1",
		bytes.NewBufferString(`{"status":"invalid"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400")
	}
}
func TestTransactionHandler_Delete_Success(t *testing.T) {
	repo := &mockTransactionRepo{
		deleteFn: func(id uint) error {
			return nil
		},
	}

	r := setupTransactionRouter(repo)

	req := httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204")
	}
}
func TestTransactionHandler_GetByID_InvalidID(t *testing.T) {
	r := setupTransactionRouter(&mockTransactionRepo{})

	req := httptest.NewRequest(http.MethodGet, "/transactions/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400")
	}
}
func TestTransactionHandler_GetAll_InvalidUserID(t *testing.T) {
	repo := &mockTransactionRepo{
		findAllFn: func(filter domain.TransactionFilter) ([]domain.Transaction, error) {
			return nil, nil
		},
	}
	r := setupTransactionRouter(repo)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions?user_id=abc",
		nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400")
	}
}
func TestTransactionHandler_Create_ServiceError(t *testing.T) {
	repo := &mockTransactionRepo{
		createFn: func(tx *domain.Transaction) error {
			return errors.New("db error")
		},
	}

	r := setupTransactionRouter(repo)

	req := httptest.NewRequest(
		http.MethodPost,
		"/transactions",
		bytes.NewBufferString(`{"user_id":1,"amount":1000}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
func TestTransactionHandler_GetAll_Detailed(t *testing.T) {
	repo := &mockTransactionRepo{
		findAllFn: func(filter domain.TransactionFilter) ([]domain.Transaction, error) {
			if filter.UserID != nil && *filter.UserID == 1 {
				return []domain.Transaction{{ID: 1}}, nil
			}
			return []domain.Transaction{}, nil
		},
	}
	r := setupTransactionRouter(repo)

	t.Run("Filter by UserID and Status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/transactions?user_id=1&status=success&page=2&limit=5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		repo.findAllFn = func(filter domain.TransactionFilter) ([]domain.Transaction, error) {
			return nil, errors.New("unexpected error")
		}
		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
func TestTransactionHandler_EdgeCases(t *testing.T) {
	repo := &mockTransactionRepo{}
	r := setupTransactionRouter(repo)

	t.Run("UpdateStatus - Not Found", func(t *testing.T) {
		repo.findByIDFn = func(id uint) (*domain.Transaction, error) {
			return nil, domain.ErrTransactionNotFound
		}
		body := `{"status":"success"}`
		req := httptest.NewRequest(http.MethodPut, "/transactions/99", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Delete - Internal Error", func(t *testing.T) {
		repo.deleteFn = func(id uint) error {
			return errors.New("disk failure")
		}
		req := httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
