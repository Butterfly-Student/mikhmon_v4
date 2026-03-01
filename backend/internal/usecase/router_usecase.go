package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
)

// RouterUseCase handles router business logic
type RouterUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Manager
}

// NewRouterUseCase creates a new router use case
func NewRouterUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
) *RouterUseCase {
	return &RouterUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
	}
}

// Create creates a new router
func (uc *RouterUseCase) Create(ctx context.Context, req *entity.RouterCreateRequest) (*entity.RouterResponse, error) {
	existing, err := uc.routerRepo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("router with name '%s' already exists", req.Name)
	}

	port := req.Port
	if port == 0 {
		port = 8728
		if req.UseSSL {
			port = 8729
		}
	}

	router := &entity.Router{
		Name:        req.Name,
		Host:        req.Host,
		Port:        port,
		Username:    req.Username,
		Password:    req.Password, // TODO: Encrypt password
		UseSSL:      req.UseSSL,
		Timeout:     req.Timeout,
		IsActive:    true,
		Description: req.Description,
	}

	if router.Timeout == 0 {
		router.Timeout = 3
	}

	if err := uc.routerRepo.Create(ctx, router); err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	return router.ToResponse(), nil
}

// GetByID retrieves a router by ID
func (uc *RouterUseCase) GetByID(ctx context.Context, id uint) (*entity.RouterResponse, error) {
	router, err := uc.routerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("router not found: %w", err)
	}
	return router.ToResponse(), nil
}

// GetAll retrieves all routers
func (uc *RouterUseCase) GetAll(ctx context.Context) ([]*entity.RouterResponse, error) {
	routers, err := uc.routerRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get routers: %w", err)
	}

	responses := make([]*entity.RouterResponse, len(routers))
	for i, router := range routers {
		responses[i] = router.ToResponse()
	}
	return responses, nil
}

// Update updates a router
func (uc *RouterUseCase) Update(ctx context.Context, id uint, req *entity.RouterUpdateRequest) (*entity.RouterResponse, error) {
	router, err := uc.routerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("router not found: %w", err)
	}

	if req.Name != "" {
		router.Name = req.Name
	}
	if req.Host != "" {
		router.Host = req.Host
	}
	if req.Port != 0 {
		router.Port = req.Port
	}
	if req.Username != "" {
		router.Username = req.Username
	}
	if req.Password != "" {
		router.Password = req.Password // TODO: Encrypt
	}
	if req.Timeout != 0 {
		router.Timeout = req.Timeout
	}
	router.UseSSL = req.UseSSL
	router.IsActive = req.IsActive
	router.Description = req.Description

	if err := uc.routerRepo.Update(ctx, router); err != nil {
		return nil, fmt.Errorf("failed to update router: %w", err)
	}

	return router.ToResponse(), nil
}

// Delete deletes a router
func (uc *RouterUseCase) Delete(ctx context.Context, id uint) error {
	if err := uc.routerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete router: %w", err)
	}
	return nil
}

// TestConnection tests connection to a stored router
func (uc *RouterUseCase) TestConnection(ctx context.Context, id uint) error {
	router, err := uc.routerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("router not found: %w", err)
	}

	timeout := time.Duration(router.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	cfg := mikrotik.Config{
		Host:     router.Host,
		Port:     router.Port,
		Username: router.Username,
		Password: router.Password,
		UseTLS:   router.UseSSL,
		Timeout:  timeout,
	}

	if err := uc.mikrotikSvc.TestConnection(ctx, cfg); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	uc.routerRepo.UpdateLastConnected(ctx, id)
	return nil
}

// TestConnectionWithParams tests a connection with the provided parameters
func (uc *RouterUseCase) TestConnectionWithParams(ctx context.Context, req *entity.RouterConnectionRequest) error {
	port := req.Port
	if port == 0 {
		port = 8728
		if req.UseSSL {
			port = 8729
		}
	}

	timeout := time.Duration(req.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	cfg := mikrotik.Config{
		Host:     req.Host,
		Port:     port,
		Username: req.Username,
		Password: req.Password,
		UseTLS:   req.UseSSL,
		Timeout:  timeout,
	}

	if err := uc.mikrotikSvc.TestConnection(ctx, cfg); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	return nil
}
