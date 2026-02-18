package services

import (
	"context"

	"github.com/evythrossell/account-management-api/internal/core/port"
)

type healthService struct {
	checker port.DBHealthChecker
}

func NewHealthService(checker port.DBHealthChecker) port.HealthService {
	return &healthService{checker: checker}
}

func (service *healthService) Check(ctx context.Context) error {
	return service.checker.PingContext(ctx)
}
