package mikrotik

import (
	"context"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
)

// connectRouter lazily connects to a router and returns its Client.
// If a Client for the router is already registered in the Manager it is returned directly.
func connectRouter(ctx context.Context, mgr *mikrotik.Manager, router *entity.Router) (*mikrotik.Client, error) {
	cfg := mikrotik.Config{
		Host:     router.Host,
		Port:     router.Port,
		Username: router.Username,
		Password: router.Password,
		UseTLS:   router.UseSSL,
		Timeout:  time.Duration(router.Timeout) * time.Second,
	}
	return mgr.GetOrConnect(ctx, router.Name, cfg)
}
