package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ClearUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "user_session")
		session.Clear()
		session.Save()
		c.JSON(200, gin.H{"success": 1})
	}
}
