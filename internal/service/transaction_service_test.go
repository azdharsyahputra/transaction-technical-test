package service

import (
	"errors"
	"testing"
	"transaction-technical-test/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_Create(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewTransactionService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything).Return(nil).Once()

		tx, err := svc.Create(1, 10000)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), tx.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything).Return(errors.New("db error")).Once()

		tx, err := svc.Create(1, 10000)

		assert.Error(t, err)
		assert.Nil(t, tx)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_UpdateStatus(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewTransactionService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		tx := &domain.Transaction{ID: 1, Status: domain.StatusPending}
		mockRepo.On("FindByID", uint(1)).Return(tx, nil).Once()
		mockRepo.On("Update", mock.Anything).Return(nil).Once()

		err := svc.UpdateStatus(1, domain.StatusSuccess)
		assert.NoError(t, err)
	})

	t.Run("Invalid Status", func(t *testing.T) {
		tx := &domain.Transaction{ID: 1, Status: domain.StatusPending}
		mockRepo.On("FindByID", uint(1)).Return(tx, nil).Once()

		err := svc.UpdateStatus(1, domain.TransactionStatus("invalid"))
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidStatus))
	})
}
func TestTransactionService_Others(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewTransactionService(mockRepo)

	t.Run("GetByID - Success", func(t *testing.T) {
		mockRepo.On("FindByID", uint(1)).Return(&domain.Transaction{ID: 1}, nil).Once()
		res, err := svc.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("GetAll - Success", func(t *testing.T) {
		filter := domain.TransactionFilter{}
		mockRepo.On("FindAll", filter).Return([]domain.Transaction{}, nil).Once()
		res, err := svc.GetAll(filter)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("UpdateStatus - FindByID Error", func(t *testing.T) {
		mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found")).Once()
		err := svc.UpdateStatus(1, domain.StatusSuccess)
		assert.Error(t, err)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		mockRepo.On("Delete", uint(1)).Return(nil).Once()
		err := svc.Delete(1)
		assert.NoError(t, err)
	})
}
