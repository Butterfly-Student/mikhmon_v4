package main

import (
	"context"
	"embed"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
	"mikhmon_v4/internal/handler"
	"mikhmon_v4/internal/middleware"
	ros "mikhmon_v4/internal/routeros"
)

//go:embed web/static
var staticFiles embed.FS

//go:embed web/templates
var templateFiles embed.FS

//go:embed voucher_templates
var voucherTemplates embed.FS

func main() {
	// Load config.
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	// Session store.
	store := sessions.NewCookieStore([]byte(cfg.Server.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}

	// RouterOS connection pool.
	pool := ros.NewPool()
	defer pool.CloseAll()

	// Parse templates.
	htmlRenderer, err := loadTemplates()
	if err != nil {
		slog.Error("failed to load templates", "err", err)
		os.Exit(1)
	}

	// Gin setup.
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.HTMLRender = htmlRenderer

	// Static files.
	staticFS, _ := fs.Sub(staticFiles, "web/static")
	r.StaticFS("/static", http.FS(staticFS))

	// Favicon.
	r.GET("/favicon.png", func(c *gin.Context) {
		c.FileFromFS("web/static/img/favicon.png", http.FS(staticFiles))
	})

	// Handlers.
	authH := handler.NewAuthHandler(store)
	adminH := handler.NewAdminHandler(store)
	dashH := handler.NewDashboardHandler(store, pool)
	hotspotH := handler.NewHotspotHandler(store, pool)
	profileH := handler.NewProfileHandler(store, pool)
	voucherH := handler.NewVoucherHandler(store, pool, voucherTemplates)
	reportH := handler.NewReportHandler(store, pool)
	trafficH := handler.NewTrafficHandler(store, pool)
	systemH := handler.NewSystemHandler(store, pool)

	// ── Public routes ─────────────────────────────────────────────
	r.GET("/login", authH.LoginPage)
	r.POST("/login", authH.Login)
	r.POST("/logout", authH.Logout)
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/login") })

	// ── Authenticated routes ───────────────────────────────────────
	authed := r.Group("/")
	authed.Use(middleware.Auth(store))

	// Admin.
	authed.GET("/admin/settings", adminH.SettingsPage)
	authed.POST("/admin/api/router", adminH.AddRouter)
	authed.PUT("/admin/api/router/:name", adminH.SaveRouter)
	authed.DELETE("/admin/api/router/:name", adminH.RemoveRouter)
	authed.POST("/admin/api/template", adminH.SaveTemplate)
	authed.GET("/admin/api/template", adminH.GetTemplate)
	authed.GET("/admin/template-editor", adminH.EditorPage)
	authed.POST("/admin/api/admin", adminH.SaveAdmin)
	authed.POST("/admin/api/router/test", adminH.TestConnection)

	// Session-scoped routes.
	sess := authed.Group("/:session")

	// Dashboard.
	sess.GET("/dashboard", dashH.DashboardPage)
	sess.GET("/api/dashboard", dashH.GetDashboard)
	sess.GET("/api/connect", dashH.TestConnection)

	// Hotspot users.
	sess.GET("/hotspot", hotspotH.HotspotPage)
	sess.GET("/api/users", hotspotH.GetUsers)
	sess.GET("/api/user/:id", hotspotH.GetUser)
	sess.POST("/api/user", hotspotH.AddUser)
	sess.PUT("/api/user/:id", hotspotH.UpdateUser)
	sess.DELETE("/api/user/:id", hotspotH.RemoveUser)

	// Active sessions.
	sess.GET("/api/active", hotspotH.GetActiveSessions)
	sess.DELETE("/api/active/:id", hotspotH.RemoveActiveSession)

	// Hosts.
	sess.GET("/api/hosts", hotspotH.GetHosts)
	sess.DELETE("/api/host/:id", hotspotH.RemoveHost)

	// Hotspot servers.
	sess.GET("/api/servers", hotspotH.GetHotspotServers)

	// Profiles.
	sess.GET("/api/profiles", profileH.GetProfiles)
	sess.GET("/api/profile/:id", profileH.GetProfile)
	sess.POST("/api/profile", profileH.AddProfile)
	sess.PUT("/api/profile/:id", profileH.UpdateProfile)
	sess.DELETE("/api/profile/:id", profileH.RemoveProfile)

	// Supporting data for profiles.
	sess.GET("/api/pools", profileH.GetAddressPools)
	sess.GET("/api/queues", profileH.GetQueues)
	sess.GET("/api/nat", profileH.GetNATRules)
	sess.GET("/api/interfaces", profileH.GetInterfaces)

	// Vouchers.
	sess.GET("/generate", voucherH.GeneratePage)
	sess.POST("/api/generate", voucherH.GenerateVouchers)
	sess.POST("/api/cache-voucher", voucherH.CacheVouchers)
	sess.GET("/print-voucher", voucherH.PrintVouchers)

	// Reports.
	sess.GET("/report", reportH.ReportPage)
	sess.GET("/api/report", reportH.GetReports)
	sess.GET("/live-report", reportH.LiveReportPage)
	sess.GET("/api/livereport", reportH.GetLiveReport)

	// Traffic (SSE).
	sess.GET("/api/traffic", trafficH.GetTraffic)

	// System.
	sess.GET("/log", systemH.LogPage)
	sess.GET("/api/logs", systemH.GetLogs)
	sess.GET("/api/expire-monitor", systemH.GetExpireMonitor)
	sess.POST("/api/expire-monitor", systemH.SetExpireMonitor)

	// ── Start server ───────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		slog.Info("Server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

type MultiTemplate struct {
	templates map[string]*template.Template
}

func (m MultiTemplate) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: m.templates[name],
		Name:     name, // Render the page template directly, which then calls base
		Data:     data,
	}
}

func loadTemplates() (render.HTMLRender, error) {
	funcMap := template.FuncMap{
		"unescapeHTML": func(s string) template.HTML {
			return template.HTML(html.UnescapeString(s))
		},
	}

	// Load layout files first
	layouts := []string{
		"web/templates/layout/base.html",
		"web/templates/layout/menu.html",
		"web/templates/layout/footer.html",
	}

	m := MultiTemplate{
		templates: make(map[string]*template.Template),
	}

	err := fs.WalkDir(templateFiles, "web/templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		// Skip layout files in the walk (they will be added to each page)
		for _, l := range layouts {
			if path == l {
				return nil
			}
		}

		// Create a new template set for this specific page
		name := path
		tmpl := template.New(name).Funcs(funcMap)

		// Add layouts
		for _, l := range layouts {
			data, err := fs.ReadFile(templateFiles, l)
			if err != nil {
				return err
			}
			_, err = tmpl.New(l).Parse(string(data))
			if err != nil {
				return err
			}
		}

		// Add the page itself
		data, err := fs.ReadFile(templateFiles, path)
		if err != nil {
			return err
		}
		_, err = tmpl.Parse(string(data))
		if err != nil {
			return err
		}

		m.templates[name] = tmpl
		return nil
	})

	return m, err
}
