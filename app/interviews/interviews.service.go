package interviews

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PAGE_SIZE = 3

type CreateInterview struct {
	Description string
}

type CreateComment struct {
	InterviewId primitive.ObjectID
	Content     string
}

func GetInterviewWithComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		interviewId := c.Param("interviewId")

		interview, err := getInterviewById(c, interviewId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Get Interview With Comment": err.Error()})
			return
		}
		comments, err := getCommentByInterviewId(c, interviewId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Get Interview With Comment": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"interview": interview, "comments": comments})
	}
}

func GetInterviewsByPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Param("page")
		pageNumber, err := utils.StringToPositiveInt(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Get Interview By Page": err.Error()})
			return
		}
		interviews, err := getInterviewsPagination(c, PAGE_SIZE, pageNumber)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Get Interview By Page": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, interviews)
	}
}

func CreateNewInterview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var interview Interview
		if err := c.ShouldBindJSON(&interview); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create Interview": err.Error()})
			return
		}

		insertErr := insertInterview(c, interview)
		if insertErr != nil {
			log.Println(insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"Create Interview": insertErr.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Create Interview": "Success"})
	}
}

func CreateNewComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment CreateComment
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create Comment": err.Error()})
			return
		}

		insertErr := insertComment(c, comment)
		if insertErr != nil {
			log.Println(insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"Create Comment": insertErr.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Create Comment": "Success"})
	}
}
