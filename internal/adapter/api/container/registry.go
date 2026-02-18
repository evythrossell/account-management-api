package container

import (
	"database/sql"

	"github.com/evythrossell/account-management-api/config"
	httpadapter "github.com/evythrossell/account-management-api/internal/adapter/http/handler"
	logger "github.com/evythrossell/account-management-api/internal/adapter/logger"
	dbadapter "github.com/evythrossell/account-management-api/internal/adapter/repository"
	"github.com/evythrossell/account-management-api/internal/core/ports"
	services "github.com/evythrossell/account-management-api/internal/core/services"
	_ "github.com/lib/pq"
)

type Container struct {
	logger                logger.Logger
	db                    *sql.DB
	accountRepository     ports.AccountRepository
	transactionRepository ports.TransactionRepository
	operationRepository   ports.OperationRepository
	accountService        ports.AccountService
	transactionService    ports.TransactionService
	healthService         ports.HealthService
	accountHandler        *httpadapter.AccountHandler
	healthHandler         *httpadapter.HealthHandler
	transactionHandler    *httpadapter.TransactionHandler
}

func New(cfg *config.Config, logger logger.Logger) (*Container, error) {
	c := &Container{
		logger: logger,
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	c.db = db
	c.logger.Info("database initialized")

	c.accountRepository = dbadapter.NewPostgresAccountRepository(db)
	c.transactionRepository = dbadapter.NewPostgresTransactionRepository(db)
	c.operationRepository = dbadapter.NewPostgresOperationRepository(db)
	c.logger.Info("repositories initialized")

	c.accountService = services.NewAccountService(c.accountRepository)
	c.transactionService = services.NewTransactionService(
		c.accountRepository,
		c.transactionRepository,
		c.operationRepository,
	)
	c.healthService = services.NewHealthService(c.DB())
	c.logger.Info("services initialized")

	c.accountHandler = httpadapter.NewAccountHandler(c.accountService)
	c.transactionHandler = httpadapter.NewTransactionHandler(c.transactionService)
	c.healthHandler = httpadapter.NewHealthHandler(c.healthService)
	c.logger.Info("handlers initialized")

	return c, nil
}

func (c *Container) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Container) Logger() logger.Logger {
	return c.logger
}

func (c *Container) DB() *sql.DB {
	return c.db
}

func (c *Container) AccountRepository() ports.AccountRepository {
	return c.accountRepository
}

func (c *Container) TransactionRepository() ports.TransactionService {
	return c.transactionService
}

func (c *Container) OperationRepository() ports.OperationRepository {
	return c.operationRepository
}

func (c *Container) HealthService() ports.HealthService {
	return c.healthService
}

func (c *Container) AccountService() ports.AccountService {
	return c.accountService
}

func (c *Container) TransactionService() ports.TransactionService {
	return c.transactionService
}

func (c *Container) AccountHandler() *httpadapter.AccountHandler {
	return c.accountHandler
}

func (c *Container) HealthHandler() *httpadapter.HealthHandler {
	return c.healthHandler
}

func (c *Container) TransactionHandler() *httpadapter.TransactionHandler {
	return c.transactionHandler
}
