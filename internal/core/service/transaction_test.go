package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	services "github.com/evythrossell/account-management-api/internal/core/service"
	common "github.com/evythrossell/account-management-api/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct{ mock.Mock }

func (m *MockTransactionRepository) Save(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	args := m.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByTransactionID(ctx context.Context, id int64) (*domain.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

type MockOperationRepository struct{ mock.Mock }

func (m *MockOperationRepository) Exists(ctx context.Context, id int16) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestTransactionService(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateTransaction - Success", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		txRepo := new(MockTransactionRepository)
		opRepo := new(MockOperationRepository)
		svc := services.NewTransactionService(accRepo, txRepo, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{ID: 1}, nil)
		opRepo.On("Exists", ctx, int16(4)).Return(true, nil)
		txRepo.On("Save", ctx, mock.Anything).Return(&domain.Transaction{ID: 100}, nil)

		res, err := svc.CreateTransaction(ctx, 1, 4, 50.0)

		assert.NoError(t, err)
		assert.Equal(t, int64(100), res.ID)
	})

	t.Run("CreateTransaction - Account Not Found", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		svc := services.NewTransactionService(accRepo, nil, nil)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(nil, common.ErrAccountNotFound)

		_, err := svc.CreateTransaction(ctx, 1, 4, 50.0)

		assert.ErrorIs(t, err, common.ErrAccountNotFound)
	})

	t.Run("CreateTransaction - Account Error (other than not found)", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		svc := services.NewTransactionService(accRepo, nil, nil)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(nil, errors.New("db connection error"))

		_, err := svc.CreateTransaction(ctx, 1, 4, 50.0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("CreateTransaction - OpRepo Error", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		svc := services.NewTransactionService(accRepo, nil, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{}, nil)
		opRepo.On("Exists", ctx, mock.Anything).Return(false, errors.New("db error"))

		_, err := svc.CreateTransaction(ctx, 1, 4, 50.0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("CreateTransaction - Invalid Operation Type", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		svc := services.NewTransactionService(accRepo, nil, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{}, nil)
		opRepo.On("Exists", ctx, mock.Anything).Return(false, nil)

		_, err := svc.CreateTransaction(ctx, 1, 99, 50.0)

		assert.ErrorIs(t, err, common.ErrInvalidOperation)
	})

	t.Run("CreateTransaction - Domain Validation Error", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		svc := services.NewTransactionService(accRepo, nil, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{}, nil)
		opRepo.On("Exists", ctx, mock.Anything).Return(true, nil)

		_, err := svc.CreateTransaction(ctx, 1, 4, -10.0)

		assert.ErrorIs(t, err, common.ErrInvalidAmount)
	})

	t.Run("CreateTransaction - Save Error", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		txRepo := new(MockTransactionRepository)
		svc := services.NewTransactionService(accRepo, txRepo, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{}, nil)
		opRepo.On("Exists", ctx, mock.Anything).Return(true, nil)
		txRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("save error"))

		_, err := svc.CreateTransaction(ctx, 1, 4, 50.0)

		assert.EqualError(t, err, "save error")
	})

	t.Run("CreateTransaction - Payment Operation", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		txRepo := new(MockTransactionRepository)
		svc := services.NewTransactionService(accRepo, txRepo, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{ID: 1}, nil)
		opRepo.On("Exists", ctx, int16(1)).Return(true, nil) // Payment is 1
		txRepo.On("Save", ctx, mock.Anything).Return(&domain.Transaction{ID: 50}, nil)

		res, err := svc.CreateTransaction(ctx, 1, 1, 25.0)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("CreateTransaction - Domain Error (generic)", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		svc := services.NewTransactionService(accRepo, nil, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{}, nil)
		opRepo.On("Exists", ctx, mock.Anything).Return(true, nil)

		_, err := svc.CreateTransaction(ctx, 1, 99, 50.0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidOperation)
	})

	t.Run("CreateTransaction - Withdrawal Operation", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		txRepo := new(MockTransactionRepository)
		svc := services.NewTransactionService(accRepo, txRepo, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{ID: 1}, nil)
		opRepo.On("Exists", ctx, int16(3)).Return(true, nil) // Withdrawal is 3
		txRepo.On("Save", ctx, mock.Anything).Return(&domain.Transaction{ID: 51}, nil)

		res, err := svc.CreateTransaction(ctx, 1, 3, 15.0)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("GetByTransactionID - Success", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		svc := services.NewTransactionService(nil, txRepo, nil)

		txRepo.On("FindByTransactionID", ctx, int64(100)).Return(&domain.Transaction{ID: 100}, nil)

		res, err := svc.GetByTransactionID(ctx, 100)

		assert.NoError(t, err)
		assert.Equal(t, int64(100), res.ID)
	})

	t.Run("GetByTransactionID - Not Found", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		svc := services.NewTransactionService(nil, txRepo, nil)

		txRepo.On("FindByTransactionID", ctx, int64(999)).Return(nil, common.ErrTransactionNotFound)

		res, err := svc.GetByTransactionID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, common.ErrTransactionNotFound)
	})

	t.Run("GetByTransactionID - Database Error", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		svc := services.NewTransactionService(nil, txRepo, nil)

		txRepo.On("FindByTransactionID", ctx, int64(100)).Return(nil, errors.New("connection failed"))

		res, err := svc.GetByTransactionID(ctx, 100)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("GetByTransactionID - Generic Error", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		svc := services.NewTransactionService(nil, txRepo, nil)

		txRepo.On("FindByTransactionID", ctx, int64(100)).Return(nil, errors.New("not found"))

		_, err := svc.GetByTransactionID(ctx, 100)

		assert.Error(t, err)
	})

	t.Run("CreateTransaction - OpRepo Exists Error", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		svc := services.NewTransactionService(accRepo, nil, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{}, nil)
		opRepo.On("Exists", ctx, mock.Anything).Return(false, errors.New("database error"))

		_, err := svc.CreateTransaction(ctx, 1, 4, 50.0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("CreateTransaction - Zero Amount", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		svc := services.NewTransactionService(accRepo, nil, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{}, nil)
		opRepo.On("Exists", ctx, mock.Anything).Return(true, nil)

		_, err := svc.CreateTransaction(ctx, 1, 4, 0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidAmount)
	})

	t.Run("CreateTransaction - Negative Amount", func(t *testing.T) {
		accRepo := new(MockAccountRepository)
		opRepo := new(MockOperationRepository)
		svc := services.NewTransactionService(accRepo, nil, opRepo)

		accRepo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{}, nil)
		opRepo.On("Exists", ctx, mock.Anything).Return(true, nil)

		_, err := svc.CreateTransaction(ctx, 1, 4, -50.0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidAmount)
	})
}
