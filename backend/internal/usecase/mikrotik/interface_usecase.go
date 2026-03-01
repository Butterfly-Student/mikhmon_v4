package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// InterfaceUseCase handles interface business logic
type InterfaceUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Manager
	log         *zap.Logger
}

// NewInterfaceUseCase creates a new interface use case
func NewInterfaceUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	log *zap.Logger,
) *InterfaceUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &InterfaceUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("interface-usecase"),
	}
}

// GetInterfaces retrieves network interfaces
func (uc *InterfaceUseCase) GetInterfaces(ctx context.Context, routerID uint) ([]*dto.Interface, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	c, err := connectRouter(ctx, uc.mikrotikSvc, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}

	interfacesCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetInterfaces(interfacesCtx)
}

// GetTraffic retrieves traffic stats for an interface (single reading)
func (uc *InterfaceUseCase) GetTraffic(ctx context.Context, routerID uint, iface string) (*dto.TrafficStats, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	c, err := connectRouter(ctx, uc.mikrotikSvc, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}

	monitorCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resultChan := make(chan dto.TrafficMonitorStats, 1)
	_, err = c.StartTrafficMonitorListen(monitorCtx, iface, resultChan)
	if err != nil {
		return nil, err
	}

	select {
	case stats := <-resultChan:
		return &dto.TrafficStats{
			Name:                  stats.Name,
			RxBitsPerSecond:       stats.RxBitsPerSecond,
			TxBitsPerSecond:       stats.TxBitsPerSecond,
			RxPacketsPerSecond:    stats.RxPacketsPerSecond,
			TxPacketsPerSecond:    stats.TxPacketsPerSecond,
			FpRxBitsPerSecond:     stats.FpRxBitsPerSecond,
			FpTxBitsPerSecond:     stats.FpTxBitsPerSecond,
			FpRxPacketsPerSecond:  stats.FpRxPacketsPerSecond,
			FpTxPacketsPerSecond:  stats.FpTxPacketsPerSecond,
			RxDropsPerSecond:      stats.RxDropsPerSecond,
			TxDropsPerSecond:      stats.TxDropsPerSecond,
			TxQueueDropsPerSecond: stats.TxQueueDropsPerSecond,
			RxErrorsPerSecond:     stats.RxErrorsPerSecond,
			TxErrorsPerSecond:     stats.TxErrorsPerSecond,
		}, nil
	case <-monitorCtx.Done():
		return &dto.TrafficStats{Name: iface}, nil
	}
}
