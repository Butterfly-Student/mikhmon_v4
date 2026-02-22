package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

const csrfKey = "csrf_token"

// CSRF generates a token on GET requests and validates it on POST/PUT/DELETE.
func CSRF(store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, _ := store.Get(c.Request, "mikhmon")

		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			// Ensure a token exists.
			if sess.Values[csrfKey] == nil {
				sess.Values[csrfKey] = newToken()
				sess.Save(c.Request, c.Writer)
			}
			c.Set(csrfKey, sess.Values[csrfKey])
		default:
			// Validate token from header or form.
			sessionToken, _ := sess.Values[csrfKey].(string)
			requestToken := c.GetHeader("X-CSRF-Token")
			if requestToken == "" {
				requestToken = c.PostForm("_csrf")
			}
			if sessionToken == "" || sessionToken != requestToken {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid CSRF token"})
				return
			}
		}
		c.Next()
	}
}

// CSRFToken retrieves the current CSRF token from the Gin context (set by the middleware).
func CSRFToken(c *gin.Context) string {
	v, _ := c.Get(csrfKey)
	s, _ := v.(string)
	return s
}

func newToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
