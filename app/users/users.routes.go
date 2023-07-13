package users

import (
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/configs"
	"github.com/sorathank/robinhood-assignment/app/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func UsersRoutes(route *gin.Engine, db *mongo.Database, cf configs.Configuration) {
	sessionManager := middleware.NewSessionManager()
	controller := NewUserController(db, sessionManager)

	login := route.Group("/login")
	{
		login.POST("", controller.ValidateUser())
	}

	user := route.Group("/user")
	{
		user.POST("", controller.CreateNewUser())
	}
}
