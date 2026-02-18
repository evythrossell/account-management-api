package port

import "context"

type HealthService interface {
	Check(ctx context.Context) error
}

type DBHealthChecker interface {
	PingContext(ctx context.Context) error
}
