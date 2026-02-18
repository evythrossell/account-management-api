package services

import (
	"context"

	"github.com/evythrossell/account-management-api/internal/core/ports"
)

type healthService struct {
	checker ports.DBHealthChecker
}

func NewHealthService(checker ports.DBHealthChecker) ports.HealthService {
	return &healthService{checker: checker}
}

func (service *healthService) Check(ctx context.Context) error {
	return service.checker.PingContext(ctx)
}
