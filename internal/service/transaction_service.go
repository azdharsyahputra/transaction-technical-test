package service

import "transaction-technical-test/internal/domain"

type TransactionService struct {
	repo domain.TransactionRepository
}

func NewTransactionService(repo domain.TransactionRepository) *TransactionService {
	return &TransactionService{
		repo: repo,
	}
}

// Create transaksi baru
func (s *TransactionService) Create(userID uint, amount float64) (*domain.Transaction, error) {
	tx := domain.NewTransaction(userID, amount)

	if err := s.repo.Create(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// GetByID ambil transaksi berdasarkan ID
func (s *TransactionService) GetByID(id uint) (*domain.Transaction, error) {
	return s.repo.FindByID(id)
}

// GetAll ambil list transaksi dengan filter
func (s *TransactionService) GetAll(filter domain.TransactionFilter) ([]domain.Transaction, error) {
	return s.repo.FindAll(filter)
}

// UpdateStatus update status transaksi
func (s *TransactionService) UpdateStatus(id uint, status domain.TransactionStatus) error {
	tx, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if err := tx.UpdateStatus(status); err != nil {
		return err
	}

	return s.repo.Update(tx)
}

// Delete hapus transaksi
func (s *TransactionService) Delete(id uint) error {
	return s.repo.Delete(id)
}
