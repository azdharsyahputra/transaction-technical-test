package service

import "transaction-technical-test/internal/domain"

type DashboardSummary struct {
	TotalSuccessToday    float64              `json:"total_success_today"`
	AverageAmountPerUser float64              `json:"average_amount_per_user"`
	LatestTransactions   []domain.Transaction `json:"latest_transactions"`
}

type DashboardService struct {
	repo domain.TransactionRepository
}

func NewDashboardService(repo domain.TransactionRepository) *DashboardService {
	return &DashboardService{
		repo: repo,
	}
}

func (s *DashboardService) GetSummary() (*DashboardSummary, error) {
	totalToday, err := s.repo.TotalSuccessToday()
	if err != nil {
		return nil, err
	}

	avgPerUser, err := s.repo.AverageAmountPerUser()
	if err != nil {
		return nil, err
	}

	latest, err := s.repo.Latest(10)
	if err != nil {
		return nil, err
	}

	return &DashboardSummary{
		TotalSuccessToday:    totalToday,
		AverageAmountPerUser: avgPerUser,
		LatestTransactions:   latest,
	}, nil
}
