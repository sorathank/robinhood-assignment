package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "user_session")
		username := session.Get("username")
		if username == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
		}
	}
}
