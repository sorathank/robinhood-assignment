package interviews

import (
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/configs"
)

type InterviewController struct {
	CF configs.Configuration
}

func InterviewRoutes(route *gin.Engine, cf configs.Configuration) {
	controller := InterviewController{cf}
	interview := route.Group("/interview")
	{
		interview.POST("", controller.CreateNewInterview())
		interview.GET("/id/:interviewId", controller.GetInterviewWithComment())
		interview.GET("/page/:page", controller.GetInterviewsByPage())
		interview.PUT("/status", controller.UpdateInterviewStatus())
	}

	comment := route.Group("/comment")
	{
		comment.POST("", controller.CreateNewComment())
	}

}
