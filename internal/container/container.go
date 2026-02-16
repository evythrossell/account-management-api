package container

import (
	"database/sql"

	"github.com/evythrossell/account-management-api/config"
	dbadapter "github.com/evythrossell/account-management-api/internal/adapters/db"
	httpadapter "github.com/evythrossell/account-management-api/internal/adapters/http"
	logger "github.com/evythrossell/account-management-api/internal/adapters/logger"
	"github.com/evythrossell/account-management-api/internal/core/ports"
	services "github.com/evythrossell/account-management-api/internal/core/services"
	_ "github.com/lib/pq"
)

type Container struct {
	logger            logger.Logger
	db                *sql.DB
	accountRepository ports.AccountRepository
	accountService    ports.AccountService
	accountHandler    *httpadapter.AccountHandler
	healthHandler     *httpadapter.HealthHandler
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

	c.accountRepository = dbadapter.NewAccountRepository(db)
	c.logger.Info("repositories initialized")

	c.accountService = services.NewAccountService(c.accountRepository)
	c.logger.Info("services initialized")

	c.accountHandler = httpadapter.NewAccountHandler(c.accountService)
	c.logger.Info("handlers initialized")

	c.healthHandler = httpadapter.NewHealthHandler(c.db)
	c.logger.Info("health handler initialized")

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

func (c *Container) AccountService() ports.AccountService {
	return c.accountService
}

func (c *Container) AccountHandler() *httpadapter.AccountHandler {
	return c.accountHandler
}

func (c *Container) HealthHandler() *httpadapter.HealthHandler {
	return c.healthHandler
}
