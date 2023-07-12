package session

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := GetUserSessionInternal(c)

		c.JSON(200, gin.H{"username": response})
		fmt.Println(response)
	}
}

func GetUserSessionInternal(c *gin.Context) interface{} {
	session := sessions.DefaultMany(c, "user_session")
	return session.Get("username")
}

func clearUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "user_session")
		session.Clear()
		session.Save()
		c.JSON(200, gin.H{"success": 1})
	}
}
