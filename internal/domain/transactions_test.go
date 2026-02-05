package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	userID := uint(1)
	amount := 50000.0

	tx := NewTransaction(userID, amount)

	assert.NotNil(t, tx)
	assert.Equal(t, userID, tx.UserID)
	assert.Equal(t, amount, tx.Amount)
	assert.Equal(t, StatusPending, tx.Status)
	assert.NotZero(t, tx.CreatedAt)
}

func TestUpdateStatus(t *testing.T) {
	// Table-driven test: ngetes banyak case sekaligus
	tests := []struct {
		name          string
		initialStatus TransactionStatus
		updateStatus  TransactionStatus
		wantErr       error
	}{
		{
			name:          "Update to Success - Valid",
			initialStatus: StatusPending,
			updateStatus:  StatusSuccess,
			wantErr:       nil,
		},
		{
			name:          "Update to Failed - Valid",
			initialStatus: StatusPending,
			updateStatus:  StatusFailed,
			wantErr:       nil,
		},
		{
			name:          "Update to Invalid Status - Error",
			initialStatus: StatusPending,
			updateStatus:  TransactionStatus("iseng"),
			wantErr:       ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{Status: tt.initialStatus}
			err := tx.UpdateStatus(tt.updateStatus)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.updateStatus, tx.Status)
			}
		})
	}
}
