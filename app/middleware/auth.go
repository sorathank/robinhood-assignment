package middleware

import (
	"log"
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
	CreateUserSession(c *gin.Context, username string)
}

type sessionManager struct {
}

func (s *sessionManager) GetCurrentUsername(c *gin.Context) (string, bool) {
	session := sessions.DefaultMany(c, "user_session")
	username, exists := session.Get("username").(string)
	return username, exists
}

func (s *sessionManager) CreateUserSession(c *gin.Context, username string) {
	session := sessions.DefaultMany(c, "user_session")
	session.Set("username", username)
	err := session.Save()
	if err != nil {
		log.Printf("Error saving session: %v", err)
	}
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
