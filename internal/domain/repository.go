package domain

import "time"

// TransactionFilter untuk query list transaksi
type TransactionFilter struct {
	UserID *uint
	Status *TransactionStatus
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}

// TransactionRepository adalah kontrak repository
type TransactionRepository interface {
	Create(tx *Transaction) error
	FindByID(id uint) (*Transaction, error)
	FindAll(filter TransactionFilter) ([]Transaction, error)
	Update(tx *Transaction) error
	Delete(id uint) error

	// Dashboard queries
	TotalSuccessToday() (float64, error)
	AverageAmountPerUser() (float64, error)
	Latest(limit int) ([]Transaction, error)
}
