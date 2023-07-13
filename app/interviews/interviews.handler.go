package interviews

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/middleware"
	"github.com/sorathank/robinhood-assignment/app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

const PAGE_SIZE = 3

type InterviewController struct {
	InterviewRepo  *InterviewRepository
	CommentRepo    *CommentRepository
	sessionManager middleware.SessionManager
}

type CreateComment struct {
	InterviewId string
	Content     string
}

type UpdateStatus struct {
	InterviewId string `binding:"required"`
	Status      string `binding:"required"`
}

func NewInterviewController(db *mongo.Database, sessionManager middleware.SessionManager) *InterviewController {
	return &InterviewController{
		InterviewRepo:  NewInterviewRepository(db),
		CommentRepo:    NewCommentRepository(db),
		sessionManager: sessionManager,
	}
}

func (ctr *InterviewController) GetInterviewWithComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		interviewId := c.Param("interviewId")

		interview, err := ctr.InterviewRepo.FindOneByID(interviewId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Get Interview With Comment": err.Error()})
			return
		}

		comments, err := ctr.CommentRepo.FindByInterviewID(interviewId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Get Interview With Comment": err.Error()})
			return
		}

		if comments == nil {
			comments = []Comment{}
		}

		c.JSON(http.StatusOK, gin.H{"interview": interview, "comments": comments})
	}
}

func (ctr *InterviewController) GetInterviewsByPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Param("page")
		pageNumber, err := utils.StringToPositiveInt(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Get Interview By Page": err.Error()})
			return
		}
		interviews, err := ctr.InterviewRepo.FindWithPagination(PAGE_SIZE, pageNumber)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Get Interview By Page": err.Error()})
			return
		}

		if interviews == nil {
			interviews = []Interview{}
		}

		c.JSON(http.StatusOK, interviews)
	}
}

func (ctr *InterviewController) CreateNewInterview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var interview Interview
		if err := c.ShouldBindJSON(&interview); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create Interview": err.Error()})
			return
		}

		username, exist := ctr.sessionManager.GetCurrentUsername(c)

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"Create Interview": "Not Sign in"})
			return
		}

		err := ctr.InterviewRepo.Insert(interview, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Create Interview": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Create Interview": "Success"})
	}
}

func (ctr *InterviewController) CreateNewComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment CreateComment
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create Comment": "Invalid Body"})
			return
		}

		username, exist := ctr.sessionManager.GetCurrentUsername(c)

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"Create Interview": "Not Sign in"})
			return
		}

		interview, err := ctr.InterviewRepo.FindOneByID(comment.InterviewId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create Comment": err.Error()})
			return
		}

		err = ctr.CommentRepo.Insert(Comment{
			InterviewId: interview.Id,
			Content:     comment.Content,
		}, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Create Comment": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Create Comment": "Success"})
	}
}

func IsValidStatus(status string) (Status, error) {
	switch strings.ToLower(status) {
	case "todo":
		return Todo, nil
	case "in progress":
		return InProgress, nil
	case "done":
		return Done, nil
	case "archived":
		return Archived, nil
	default:
		return "", errors.New("invalid status")
	}
}

func (ctr *InterviewController) UpdateInterviewStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateStatus UpdateStatus
		if err := c.ShouldBindJSON(&updateStatus); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Update Status": err.Error()})
			return
		}

		validStatus, err := IsValidStatus(updateStatus.Status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Update Status": err.Error()})
			return
		}

		_, err = ctr.InterviewRepo.FindOneByID(updateStatus.InterviewId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Update Status": err.Error()})
			return
		}

		err = ctr.InterviewRepo.UpdateStatus(validStatus, updateStatus.InterviewId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Update Status": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"Update Status": "Success"})
	}
}
