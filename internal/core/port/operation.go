package port

import "context"

type OperationRepository interface {
	Exists(ctx context.Context, operationType int16) (bool, error)
}
