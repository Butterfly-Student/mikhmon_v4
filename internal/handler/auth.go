package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
)

// AuthHandler handles login/logout.
type AuthHandler struct {
	store sessions.Store
}

func NewAuthHandler(store sessions.Store) *AuthHandler {
	return &AuthHandler{store: store}
}

// LoginPage renders the login form.
func (h *AuthHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "web/templates/login.html", gin.H{
		"Title": "Mikhmon — Login",
	})
}

// Login processes a POST form with username/password.
func (h *AuthHandler) Login(c *gin.Context) {
	u := c.PostForm("username")
	p := c.PostForm("password")

	cfg := config.Get()
	if cfg == nil {
		c.HTML(http.StatusInternalServerError, "web/templates/login.html", gin.H{"Error": "Config not loaded"})
		return
	}

	if u != cfg.Admin.Username || !config.CheckPassword(p, cfg.Admin.PasswordHash) {
		c.HTML(http.StatusUnauthorized, "web/templates/login.html", gin.H{"Error": "Invalid username or password"})
		return
	}

	sess, _ := h.store.Get(c.Request, "mikhmon")
	sess.Values["mikhmon"] = true
	sess.Values["username"] = u
	sess.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/admin/settings")
}

// Logout destroys the session and redirects to login.
func (h *AuthHandler) Logout(c *gin.Context) {
	sess, _ := h.store.Get(c.Request, "mikhmon")
	sess.Options.MaxAge = -1
	sess.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/login")
}
