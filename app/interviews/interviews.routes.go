package interviews

import (
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/configs"
	"github.com/sorathank/robinhood-assignment/app/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func InterviewRoutes(route *gin.Engine, db *mongo.Database, cf configs.Configuration, auth middleware.AuthMiddleware) {
	sessionManager := middleware.NewSessionManager()
	controller := NewInterviewController(db, sessionManager)
	interviewAndComment := route.Group("/")
	{
		interviewAndComment.POST("/interview", auth.RequireAuthentication(), controller.CreateNewInterview())
		interviewAndComment.GET("/interview/id/:interviewId", controller.GetInterviewWithComment())
		interviewAndComment.GET("/interview/page/:page", controller.GetInterviewsByPage())
		interviewAndComment.PUT("/interview/status", auth.RequireAuthentication(), controller.UpdateInterviewStatus())
		interviewAndComment.POST("/comment", auth.RequireAuthentication(), controller.CreateNewComment())
	}
}
