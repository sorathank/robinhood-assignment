package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireAuthentication() gin.HandlerFunc
}

type authMiddleware struct {
	sessionManager SessionManager
}

type SessionManager interface {
	GetCurrentUsername(c *gin.Context) (string, bool)
}

type sessionManager struct {
}

func (s *sessionManager) GetCurrentUsername(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	username, exists := session.Get("username").(string)
	return username, exists
}

func NewAuthMiddleware(sessionManager SessionManager) AuthMiddleware {
	return &authMiddleware{
		sessionManager: sessionManager,
	}
}

func (a *authMiddleware) RequireAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := a.sessionManager.GetCurrentUsername(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}

func NewSessionManager() SessionManager {
	return &sessionManager{}
}
