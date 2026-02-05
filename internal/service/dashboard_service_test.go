package service

import (
	"errors"
	"testing"
	"transaction-technical-test/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestDashboardService_GetSummary(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewDashboardService(mockRepo) // variabel namanya svc

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("TotalSuccessToday").Return(50000.0, nil).Once()
		mockRepo.On("AverageAmountPerUser").Return(25000.0, nil).Once()
		mockRepo.On("Latest", 10).Return([]domain.Transaction{{ID: 1}}, nil).Once()

		summary, err := svc.GetSummary() // Pastikan panggil svc, bukan s

		assert.NoError(t, err)
		assert.NotNil(t, summary)
		assert.Equal(t, 50000.0, summary.TotalSuccessToday)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error on TotalSuccess", func(t *testing.T) {
		mockRepo.On("TotalSuccessToday").Return(0.0, errors.New("db error")).Once()

		summary, err := svc.GetSummary() // Sini juga ganti jadi svc

		assert.Error(t, err)
		assert.Nil(t, summary)
		mockRepo.AssertExpectations(t)
	})
}
func TestDashboardService_GetSummary_MoreErrors(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewDashboardService(mockRepo)

	t.Run("Error on AverageAmount", func(t *testing.T) {
		mockRepo.On("TotalSuccessToday").Return(5000.0, nil).Once()
		mockRepo.On("AverageAmountPerUser").Return(0.0, errors.New("error avg")).Once()

		res, err := svc.GetSummary()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Error on Latest", func(t *testing.T) {
		mockRepo.On("TotalSuccessToday").Return(5000.0, nil).Once()
		mockRepo.On("AverageAmountPerUser").Return(2000.0, nil).Once()
		mockRepo.On("Latest", 10).Return(nil, errors.New("error latest")).Once()

		res, err := svc.GetSummary()
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
