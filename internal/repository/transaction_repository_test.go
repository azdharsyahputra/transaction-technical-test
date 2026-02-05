package repository

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"transaction-technical-test/internal/domain"
)

func setupTestRepo(t *testing.T) *TransactionRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed open db: %v", err)
	}

	if err := db.AutoMigrate(&TransactionModel{}); err != nil {
		t.Fatalf("failed migrate: %v", err)
	}

	return NewTransactionRepository(db)
}
func TestTransactionRepository_Create(t *testing.T) {
	repo := setupTestRepo(t)

	tx := &domain.Transaction{
		UserID: 1,
		Amount: 1000,
		Status: domain.StatusPending,
	}

	if err := repo.Create(tx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tx.ID == 0 {
		t.Fatalf("expected id to be set")
	}
}
func TestTransactionRepository_FindByID_Success(t *testing.T) {
	repo := setupTestRepo(t)

	tx := &domain.Transaction{
		UserID: 1,
		Amount: 500,
		Status: domain.StatusSuccess,
	}
	_ = repo.Create(tx)

	found, err := repo.FindByID(tx.ID)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if found.ID != tx.ID {
		t.Fatalf("wrong id")
	}
}
func TestTransactionRepository_FindAll_Filter(t *testing.T) {
	repo := setupTestRepo(t)

	now := time.Now()

	_ = repo.Create(&domain.Transaction{
		UserID:    1,
		Amount:    1000,
		Status:    domain.StatusSuccess,
		CreatedAt: now,
	})
	_ = repo.Create(&domain.Transaction{
		UserID:    2,
		Amount:    2000,
		Status:    domain.StatusPending,
		CreatedAt: now,
	})

	status := domain.StatusSuccess
	filter := domain.TransactionFilter{
		Status: &status,
		Limit:  10,
		Offset: 0,
	}

	result, err := repo.FindAll(filter)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result")
	}
}
func TestTransactionRepository_Update_Success(t *testing.T) {
	repo := setupTestRepo(t)

	tx := &domain.Transaction{
		UserID: 1,
		Amount: 1000,
		Status: domain.StatusPending,
	}
	_ = repo.Create(tx)

	tx.Status = domain.StatusSuccess

	if err := repo.Update(tx); err != nil {
		t.Fatalf("unexpected error")
	}
}
func TestTransactionRepository_Update_NotFound(t *testing.T) {
	repo := setupTestRepo(t)

	tx := &domain.Transaction{
		ID:     999,
		Status: domain.StatusSuccess,
	}

	err := repo.Update(tx)
	if err != domain.ErrTransactionNotFound {
		t.Fatalf("expected not found error")
	}
}
func TestTransactionRepository_Delete_Success(t *testing.T) {
	repo := setupTestRepo(t)

	tx := &domain.Transaction{
		UserID: 1,
		Amount: 100,
		Status: domain.StatusPending,
	}
	_ = repo.Create(tx)

	if err := repo.Delete(tx.ID); err != nil {
		t.Fatalf("unexpected error")
	}
}
func TestTransactionRepository_Delete_NotFound(t *testing.T) {
	repo := setupTestRepo(t)

	err := repo.Delete(999)
	if err != domain.ErrTransactionNotFound {
		t.Fatalf("expected not found error")
	}
}
func TestTransactionRepository_TotalSuccessToday(t *testing.T) {
	repo := setupTestRepo(t)

	_ = repo.Create(&domain.Transaction{
		UserID: 1,
		Amount: 1000,
		Status: domain.StatusSuccess,
	})

	total, err := repo.TotalSuccessToday()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if total <= 0 {
		t.Fatalf("expected total > 0")
	}
}
func TestTransactionRepository_AverageAmountPerUser(t *testing.T) {
	repo := setupTestRepo(t)

	_ = repo.Create(&domain.Transaction{
		UserID: 1,
		Amount: 1000,
		Status: domain.StatusSuccess,
	})
	_ = repo.Create(&domain.Transaction{
		UserID: 1,
		Amount: 3000,
		Status: domain.StatusSuccess,
	})

	avg, err := repo.AverageAmountPerUser()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if avg <= 0 {
		t.Fatalf("expected avg > 0")
	}
}
func TestTransactionRepository_Latest(t *testing.T) {
	repo := setupTestRepo(t)

	_ = repo.Create(&domain.Transaction{
		UserID: 1,
		Amount: 1000,
		Status: domain.StatusSuccess,
	})

	result, err := repo.Latest(5)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if len(result) == 0 {
		t.Fatalf("expected result")
	}
}
