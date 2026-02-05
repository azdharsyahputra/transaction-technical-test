package domain

import "time"

// TransactionStatus merepresentasikan status transaksi
type TransactionStatus string

const (
	StatusPending TransactionStatus = "pending"
	StatusSuccess TransactionStatus = "success"
	StatusFailed  TransactionStatus = "failed"
)

// Transaction adalah entity utama domain
type Transaction struct {
	ID        uint
	UserID    uint
	Amount    float64
	Status    TransactionStatus
	CreatedAt time.Time
}

// NewTransaction adalah constructor transaksi baru
func NewTransaction(userID uint, amount float64) *Transaction {
	return &Transaction{
		UserID:    userID,
		Amount:    amount,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}

// UpdateStatus mengubah status transaksi dengan validasi
func (t *Transaction) UpdateStatus(status TransactionStatus) error {
	if !isValidStatus(status) {
		return ErrInvalidStatus
	}

	t.Status = status
	return nil
}

func isValidStatus(status TransactionStatus) bool {
	switch status {
	case StatusPending, StatusSuccess, StatusFailed:
		return true
	default:
		return false
	}
}
