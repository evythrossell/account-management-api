package container

import (
	"database/sql"

	"github.com/evythrossell/account-management-api/internal/adapter/http/handler"
	dbadapter "github.com/evythrossell/account-management-api/internal/adapter/storage/postgres"
	"github.com/evythrossell/account-management-api/internal/core/port"
	service "github.com/evythrossell/account-management-api/internal/core/service"
	config "github.com/evythrossell/account-management-api/internal/infrastructure"
	logger "github.com/evythrossell/account-management-api/pkg"
	_ "github.com/lib/pq"
)

type Container struct {
	logger                logger.Logger
	db                    *sql.DB
	accountRepository     port.AccountRepository
	transactionRepository port.TransactionRepository
	operationRepository   port.OperationRepository
	accountService        port.AccountService
	transactionService    port.TransactionService
	healthService         port.HealthService
	accountHandler        *handler.AccountHandler
	healthHandler         *handler.HealthHandler
	transactionHandler    *handler.TransactionHandler
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

	c.accountService = service.NewAccountService(c.accountRepository)
	c.transactionService = service.NewTransactionService(
		c.accountRepository,
		c.transactionRepository,
		c.operationRepository,
	)
	c.healthService = service.NewHealthService(c.DB())
	c.logger.Info("services initialized")

	c.accountHandler = handler.NewAccountHandler(c.accountService)
	c.transactionHandler = handler.NewTransactionHandler(c.transactionService)
	c.healthHandler = handler.NewHealthHandler(c.healthService)
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

func (c *Container) AccountRepository() port.AccountRepository {
	return c.accountRepository
}

func (c *Container) TransactionRepository() port.TransactionRepository {
	return c.transactionRepository
}

func (c *Container) OperationRepository() port.OperationRepository {
	return c.operationRepository
}

func (c *Container) HealthService() port.HealthService {
	return c.healthService
}

func (c *Container) AccountService() port.AccountService {
	return c.accountService
}

func (c *Container) TransactionService() port.TransactionService {
	return c.transactionService
}

func (c *Container) AccountHandler() *handler.AccountHandler {
	return c.accountHandler
}

func (c *Container) HealthHandler() *handler.HealthHandler {
	return c.healthHandler
}

func (c *Container) TransactionHandler() *handler.TransactionHandler {
	return c.transactionHandler
}
