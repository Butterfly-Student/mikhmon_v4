//go:build modern

package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appbootstrap "mikhmon_v4/internal/clean/bootstrap"
	appcfg "mikhmon_v4/internal/clean/config"
)

func main() {
	ctx := context.Background()
	cfg := appcfg.Load()

	app, err := appbootstrap.New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer app.Logger.Sync()

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "app": app.Config.AppName})
	})

	r.POST("/auth/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := app.AuthService.Login(req.Username, req.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "login success"})
	})

	app.Logger.Info("modern service started", zap.String("port", cfg.HTTPPort))
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		app.Logger.Fatal("server failed", zap.Error(err))
	}
}
