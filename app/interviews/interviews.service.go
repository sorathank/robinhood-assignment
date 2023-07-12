package interviews

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PAGE_SIZE = 3

type CreateComment struct {
	InterviewId primitive.ObjectID
	Content     string
}

type UpdateStatus struct {
	InterviewId primitive.ObjectID
	Status      Status
}

func (ctr InterviewController) GetInterviewWithComment() gin.HandlerFunc {
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

func (ctr InterviewController) GetInterviewsByPage() gin.HandlerFunc {
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

func (ctr InterviewController) CreateNewInterview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var interview Interview
		if err := c.ShouldBindJSON(&interview); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create Interview": err.Error()})
			return
		}

		session := sessions.DefaultMany(c, "user_session")

		if username := session.Get("username"); username == nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create Interview": "Not Sign in"})
			return
		}

		creator := session.Get("username").(string)
		insertObject := Interview{
			Description: interview.Description,
			User:        creator,
			Status:      Todo,
			CreatedTime: time.Now(),
		}

		insertErr := insertInterview(c, insertObject)
		if insertErr != nil {
			log.Println(insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"Create Interview": insertErr.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Create Interview": "Success"})
	}
}

func (ctr InterviewController) CreateNewComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment CreateComment
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create Comment": "Invalid Body"})
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

func (ctr InterviewController) UpdateInterviewStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateStatus UpdateStatus
		if err := c.ShouldBindJSON(&updateStatus); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Update Status": "Invalid Body"})
			return
		}

		updateErr := updateInterviewStatus(c, updateStatus.Status, updateStatus.InterviewId.String())
		if updateErr != nil {
			log.Println(updateErr)
			c.JSON(http.StatusInternalServerError, gin.H{"Update Status": updateErr.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Update Status": "Success"})
	}
}
