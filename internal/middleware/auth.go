package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Auth returns a Gin middleware that requires a valid "mikhmon" session cookie.
// If the session is missing the user is redirected to /login.
func Auth(store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := store.Get(c.Request, "mikhmon")
		if err != nil || sess.Values["mikhmon"] == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
