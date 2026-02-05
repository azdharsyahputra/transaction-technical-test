package repository

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"transaction-technical-test/internal/domain"
)

func setupTestRepo(t *testing.T) *TransactionRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed open db: %v", err)
	}

	if err := db.AutoMigrate(&TransactionModel{}); err != nil {
		t.Fatalf("failed migrate: %v", err)
	}

	db.Exec("DELETE FROM transaction_models")

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

	found, err := repo.FindByID(tx.ID)
	if err != nil {
		t.Fatalf("failed find after create")
	}

	if found.Amount != 1000 || found.Status != domain.StatusPending {
		t.Fatalf("data mismatch after create")
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
	if found.Amount != 500 {
		t.Fatalf("wrong amount")
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
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].Status != domain.StatusSuccess {
		t.Fatalf("wrong status")
	}
	if result[0].Amount != 1000 {
		t.Fatalf("wrong amount")
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
	tx.Amount = 2000

	if err := repo.Update(tx); err != nil {
		t.Fatalf("unexpected error")
	}

	updated, err := repo.FindByID(tx.ID)
	if err != nil {
		t.Fatalf("failed find after update")
	}

	if updated.Status != domain.StatusSuccess {
		t.Fatalf("status not updated")
	}
	if updated.Amount != 2000 {
		t.Fatalf("amount not updated")
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

	_, err := repo.FindByID(tx.ID)
	if err != domain.ErrTransactionNotFound {
		t.Fatalf("expected not found after delete")
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

	startOfDay := time.Now().Truncate(24 * time.Hour)

	_ = repo.Create(&domain.Transaction{
		UserID:    1,
		Amount:    1000,
		Status:    domain.StatusSuccess,
		CreatedAt: startOfDay.Add(1 * time.Hour),
	})
	_ = repo.Create(&domain.Transaction{
		UserID:    2,
		Amount:    999,
		Status:    domain.StatusSuccess,
		CreatedAt: startOfDay.Add(-24 * time.Hour),
	})

	total, err := repo.TotalSuccessToday()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if total < 999.99 || total > 1000.01 {
		t.Fatalf("expected total ~1000, got %.2f", total)
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

	if avg < 1999.99 || avg > 2000.01 {
		t.Fatalf("expected avg ~2000, got %.2f", avg)
	}

}

func TestTransactionRepository_Latest(t *testing.T) {
	repo := setupTestRepo(t)

	_ = repo.Create(&domain.Transaction{
		UserID: 1,
		Amount: 1000,
		Status: domain.StatusSuccess,
	})
	_ = repo.Create(&domain.Transaction{
		UserID: 2,
		Amount: 2000,
		Status: domain.StatusPending,
	})

	result, err := repo.Latest(1)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result")
	}

	if result[0].Amount != 2000 {
		t.Fatalf("expected latest transaction")
	}
}
