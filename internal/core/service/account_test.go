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

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Save(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	args := m.Called(ctx, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByDocument(ctx context.Context, doc string) (*domain.Account, error) {
	args := m.Called(ctx, doc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByAccountID(ctx context.Context, id int64) (*domain.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func TestAccountService(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateAccount - Success", func(t *testing.T) {
		repo := new(MockAccountRepository)
		svc := services.NewAccountService(repo)
		doc := "12345678901"
		acc := &domain.Account{DocumentNumber: doc}

		repo.On("FindByDocument", ctx, doc).Return(nil, common.ErrAccountNotFound)
		repo.On("Save", ctx, mock.Anything).Return(acc, nil)

		result, err := svc.CreateAccount(ctx, doc)

		assert.NoError(t, err)
		assert.Equal(t, doc, result.DocumentNumber)
	})

	t.Run("CreateAccount - Invalid Document", func(t *testing.T) {
		svc := services.NewAccountService(nil)
		_, err := svc.CreateAccount(ctx, "invalid")
		assert.ErrorIs(t, err, common.ErrInvalidDocument)
	})

	t.Run("CreateAccount - Already Exists", func(t *testing.T) {
		repo := new(MockAccountRepository)
		svc := services.NewAccountService(repo)
		doc := "12345678901"

		repo.On("FindByDocument", ctx, doc).Return(&domain.Account{}, nil)

		_, err := svc.CreateAccount(ctx, doc)

		assert.ErrorIs(t, err, common.ErrAccountAlreadyExists)
	})

	t.Run("CreateAccount - Repository Error on Find", func(t *testing.T) {
		repo := new(MockAccountRepository)
		svc := services.NewAccountService(repo)
		doc := "12345678901"

		repo.On("FindByDocument", ctx, doc).Return(nil, errors.New("db down"))

		_, err := svc.CreateAccount(ctx, doc)

		assert.EqualError(t, err, "db down")
	})

	t.Run("GetAccount - Success", func(t *testing.T) {
		repo := new(MockAccountRepository)
		svc := services.NewAccountService(repo)
		repo.On("FindByDocument", ctx, "123").Return(&domain.Account{ID: 1}, nil)

		res, err := svc.GetAccount(ctx, "123")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.ID)
	})

	t.Run("GetAccount - Error", func(t *testing.T) {
		repo := new(MockAccountRepository)
		svc := services.NewAccountService(repo)
		repo.On("FindByDocument", ctx, "123").Return(nil, errors.New("error"))

		_, err := svc.GetAccount(ctx, "123")
		assert.Error(t, err)
	})

	t.Run("GetAccountByID - Success", func(t *testing.T) {
		repo := new(MockAccountRepository)
		svc := services.NewAccountService(repo)
		repo.On("FindByAccountID", ctx, int64(1)).Return(&domain.Account{ID: 1}, nil)

		res, err := svc.GetAccountByID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.ID)
	})

	t.Run("GetAccountByID - Error", func(t *testing.T) {
		repo := new(MockAccountRepository)
		svc := services.NewAccountService(repo)
		repo.On("FindByAccountID", ctx, int64(1)).Return(nil, errors.New("error"))

		_, err := svc.GetAccountByID(ctx, 1)
		assert.Error(t, err)
	})
}
