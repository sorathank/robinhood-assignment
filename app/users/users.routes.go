package users

import (
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/configs"
	"go.mongodb.org/mongo-driver/mongo"
)

func UsersRoutes(route *gin.Engine, db *mongo.Database, cf configs.Configuration) {
	controller := NewUserController(db)

	login := route.Group("/login")
	{
		login.POST("", controller.ValidateUser())
	}

	user := route.Group("/user")
	{
		user.POST("", controller.CreateNewUser())
	}
}
