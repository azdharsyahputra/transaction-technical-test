package service

import (
	"transaction-technical-test/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockRepo ini adalah tiruan dari TransactionRepository
type MockRepo struct {
	mock.Mock
}

// Implementasi fungsi-fungsi interface repository
func (m *MockRepo) Create(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockRepo) FindByID(id uint) (*domain.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockRepo) FindAll(filter domain.TransactionFilter) ([]domain.Transaction, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockRepo) Update(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) TotalSuccessToday() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRepo) AverageAmountPerUser() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRepo) Latest(limit int) ([]domain.Transaction, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}
