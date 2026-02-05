package repository

import (
	"time"

	"gorm.io/gorm"

	"transaction-technical-test/internal/domain"
)

type TransactionModel struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Amount    float64
	Status    string
	CreatedAt time.Time
}

// Mapper
func toDomain(m *TransactionModel) domain.Transaction {
	return domain.Transaction{
		ID:        m.ID,
		UserID:    m.UserID,
		Amount:    m.Amount,
		Status:    domain.TransactionStatus(m.Status),
		CreatedAt: m.CreatedAt,
	}
}

func fromDomain(d *domain.Transaction) TransactionModel {
	return TransactionModel{
		ID:        d.ID,
		UserID:    d.UserID,
		Amount:    d.Amount,
		Status:    string(d.Status),
		CreatedAt: d.CreatedAt,
	}
}

type TransactionRepository struct {
	db *gorm.DB
}

// Constructor
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Implement
func (r *TransactionRepository) Create(tx *domain.Transaction) error {
	model := fromDomain(tx)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	tx.ID = model.ID
	return nil
}
func (r *TransactionRepository) FindByID(id uint) (*domain.Transaction, error) {
	var model TransactionModel

	if err := r.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTransactionNotFound
		}
		return nil, err
	}

	tx := toDomain(&model)
	return &tx, nil
}
func (r *TransactionRepository) FindAll(filter domain.TransactionFilter) ([]domain.Transaction, error) {
	var models []TransactionModel

	query := r.db.Model(&TransactionModel{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}

	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}

	if filter.To != nil {
		query = query.Where("created_at <= ?", *filter.To)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset >= 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]domain.Transaction, 0, len(models))
	for _, m := range models {
		tx := toDomain(&m)
		result = append(result, tx)
	}

	return result, nil
}

func (r *TransactionRepository) Update(tx *domain.Transaction) error {
	result := r.db.Model(&TransactionModel{}).
		Where("id = ?", tx.ID).
		Updates(map[string]interface{}{
			"status": tx.Status,
			"amount": tx.Amount,
		})

	if result.RowsAffected == 0 {
		return domain.ErrTransactionNotFound
	}

	return result.Error
}

func (r *TransactionRepository) Delete(id uint) error {
	result := r.db.Delete(&TransactionModel{}, id)

	if result.RowsAffected == 0 {
		return domain.ErrTransactionNotFound
	}

	return result.Error
}
func (r *TransactionRepository) TotalSuccessToday() (float64, error) {
	var total float64

	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	err := r.db.Model(&TransactionModel{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("status = ?", string(domain.StatusSuccess)).
		Where("created_at >= ? AND created_at < ?", start, end).
		Scan(&total).Error

	return total, err
}
func (r *TransactionRepository) AverageAmountPerUser() (float64, error) {
	var avg float64

	err := r.db.Model(&TransactionModel{}).
		Select("COALESCE(AVG(amount), 0)").
		Where("status = ?", string(domain.StatusSuccess)).
		Scan(&avg).Error

	return avg, err
}

func (r *TransactionRepository) Latest(limit int) ([]domain.Transaction, error) {
	var models []TransactionModel

	if err := r.db.
		Order("created_at desc").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]domain.Transaction, 0, len(models))
	for _, m := range models {
		tx := toDomain(&m)
		result = append(result, tx)
	}

	return result, nil
}
