package db

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresOperationRepository struct {
	db *sql.DB
}

func NewPostgresOperationRepository(db *sql.DB) *PostgresOperationRepository {
	return &PostgresOperationRepository{db: db}
}

func (p *PostgresOperationRepository) Exists(ctx context.Context, operationType int16) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM operations_types WHERE operation_type_id = $1)`

	err := p.db.QueryRowContext(ctx, query, operationType).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("infrastructure error: failed to check operation existence: %w", err)
	}

	return exists, nil
}
