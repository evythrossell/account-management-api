package db

import (
	"context"
	"database/sql"

	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
)

type operationRepository struct {
	db *sql.DB
}

func NewOperationRepository(db *sql.DB) *operationRepository {
	return &operationRepository{db: db}
}

func (r *operationRepository) Exists(ctx context.Context, operationType int16) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM operations_types WHERE operation_type_id=$1)", operationType).
		Scan(&exists)
	if err != nil {
		return false, domainerror.NewConflictError("check operation exists", err)
	}

	return exists, nil
}
