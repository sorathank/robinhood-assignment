package interviews

import (
	"github.com/gin-gonic/gin"
)

func InterviewRoutes(route *gin.Engine) {
	interview := route.Group("/interview")
	{
		interview.POST("", CreateNewInterview())
		interview.GET("/id/:interviewId", GetInterviewWithComment())
		interview.GET("/page/:page", GetInterviewsByPage())
	}

}
