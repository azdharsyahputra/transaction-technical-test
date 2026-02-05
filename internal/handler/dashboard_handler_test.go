package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	"transaction-technical-test/internal/domain"
	"transaction-technical-test/internal/handler"
	"transaction-technical-test/internal/service"
)

// Mock Error
type mockDashboardErrorRepo struct{}

func (m *mockDashboardErrorRepo) TotalSuccessToday() (float64, error) {
	return 0, errors.New("db error")
}
func (m *mockDashboardErrorRepo) AverageAmountPerUser() (float64, error) {
	return 0, nil
}
func (m *mockDashboardErrorRepo) Latest(limit int) ([]domain.Transaction, error) {
	return nil, nil
}

func (m *mockDashboardErrorRepo) Create(*domain.Transaction) error { return nil }
func (m *mockDashboardErrorRepo) FindByID(uint) (*domain.Transaction, error) {
	return nil, nil
}
func (m *mockDashboardErrorRepo) FindAll(domain.TransactionFilter) ([]domain.Transaction, error) {
	return nil, nil
}
func (m *mockDashboardErrorRepo) Update(*domain.Transaction) error { return nil }
func (m *mockDashboardErrorRepo) Delete(uint) error                { return nil }

func TestDashboardHandler_Summary_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockDashboardErrorRepo{}
	svc := service.NewDashboardService(repo)

	logger := zap.NewNop()
	h := handler.NewDashboardHandler(svc, logger)

	r := gin.New()
	r.GET("/dashboard/summary", h.Summary)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500")
	}
}

// Mock Succes
type mockDashboardSuccessRepo struct{}

func (m *mockDashboardSuccessRepo) TotalSuccessToday() (float64, error) {
	return 1000, nil
}
func (m *mockDashboardSuccessRepo) AverageAmountPerUser() (float64, error) {
	return 500, nil
}
func (m *mockDashboardSuccessRepo) Latest(limit int) ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func (m *mockDashboardSuccessRepo) Create(*domain.Transaction) error { return nil }
func (m *mockDashboardSuccessRepo) FindByID(uint) (*domain.Transaction, error) {
	return nil, nil
}
func (m *mockDashboardSuccessRepo) FindAll(domain.TransactionFilter) ([]domain.Transaction, error) {
	return nil, nil
}
func (m *mockDashboardSuccessRepo) Update(*domain.Transaction) error { return nil }
func (m *mockDashboardSuccessRepo) Delete(uint) error                { return nil }

func TestDashboardHandler_Summary_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockDashboardSuccessRepo{}
	svc := service.NewDashboardService(repo)

	logger := zap.NewNop()
	h := handler.NewDashboardHandler(svc, logger)

	r := gin.New()
	r.GET("/dashboard/summary", h.Summary)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
