package interviews

import (
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/configs"
	"go.mongodb.org/mongo-driver/mongo"
)

func InterviewRoutes(route *gin.Engine, db *mongo.Database, cf configs.Configuration) {
	controller := NewInterviewController(db)
	interviewAndComment := route.Group("/")
	{
		interviewAndComment.POST("/interview", controller.CreateNewInterview())
		interviewAndComment.GET("/interview/id/:interviewId", controller.GetInterviewWithComment())
		interviewAndComment.GET("/interview/page/:page", controller.GetInterviewsByPage())
		interviewAndComment.PUT("/interview/status", controller.UpdateInterviewStatus())
		interviewAndComment.POST("/comment", controller.CreateNewComment())
	}
}
