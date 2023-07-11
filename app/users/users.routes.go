package users

import (
	"github.com/gin-gonic/gin"
)

func UsersRoutes(route *gin.Engine) {
	login := route.Group("/login")
	{
		login.POST("", ValidateUser())
	}

	user := route.Group("/user")
	{
		user.POST("", CreateNewUser())
	}
}
