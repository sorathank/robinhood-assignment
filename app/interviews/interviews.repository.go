package interviews

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Interview struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Description string             `bson:"description"`
	User        string             `bson:"user"`
	Status      Status             `bson:"status"`
	CreatedTime time.Time          `bson:"created_time"`
}

type Comment struct {
	Id          primitive.ObjectID `bson:"_id, omitempty"`
	User        string             `bson:"user"`
	InterviewId primitive.ObjectID `bson:"interview_id"`
	Content     string             `bson:"content"`
	CreatedTime time.Time          `bson:"created_time"`
}

type Status string

const (
	Todo       Status = "Todo"
	InProgress Status = "In Progress"
	Done       Status = "Done"
	Archived   Status = "Archived"
)

func getInterviewById(c *gin.Context, interviewId string) (Interview, error) {
	db := utils.GetDB(c)
	collection := db.Collection("interviews")

	var interview Interview
	err := collection.FindOne(context.TODO(), bson.M{"_id": interviewId}).Decode(&interview)

	return interview, err
}

func getCommentByInterviewId(c *gin.Context, interviewId string) ([]Comment, error) {
	db := utils.GetDB(c)
	collection := db.Collection("interviews")
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_time", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.M{"interview_id": interviewId}, findOptions)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.TODO())

	var comments []Comment
	if err := cursor.All(context.TODO(), &comments); err != nil {
		log.Println(err)
	}

	return comments, err
}

func getInterviewsPagination(c *gin.Context, pageSize int64, pageNumber int64) ([]Interview, error) {
	db := utils.GetDB(c)
	collection := db.Collection("interviews")
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_time", Value: -1}})
	findOptions.SetSkip((pageNumber - 1) * pageSize)
	findOptions.SetLimit(pageSize)

	cursor, err := collection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.TODO())

	var interviews []Interview
	if err = cursor.All(c, &interviews); err != nil {
		log.Println(err)
	}

	return interviews, err
}

func insertInterview(c *gin.Context, interview Interview) error {
	db := utils.GetDB(c)
	collection := db.Collection("interviews")

	session := sessions.DefaultMany(c, "user_session")

	insertResult, err := collection.InsertOne(context.TODO(), bson.M{
		"Description": interview.Description,
		"Creator":     session.Get("username"),
		"Status":      Todo,
		"CreatedTime": time.Now(),
	})
	fmt.Println("Inserted new Interview: ", insertResult.InsertedID)

	return err
}

func insertComment(c *gin.Context, createComment CreateComment) error {
	db := utils.GetDB(c)
	collection := db.Collection("comments")

	session := sessions.DefaultMany(c, "user_session")

	insertResult, err := collection.InsertOne(context.TODO(), bson.M{
		"User":        session.Get("username"),
		"InterviewId": createComment.InterviewId,
		"Content":     createComment.Content,
		"CreatedTime": time.Now(),
	})
	fmt.Println("Inserted new comment: ", insertResult.InsertedID)

	return err
}

func updateInterviewStatus(c *gin.Context, status Status, interviewId string) error {
	db := utils.GetDB(c)
	collection := db.Collection("interviews")
	objectId, err := primitive.ObjectIDFromHex(interviewId)
	_, err = collection.UpdateOne(
		c,
		bson.M{"_id": objectId},
		bson.D{
			{"$set", bson.D{{"status", status}}},
		},
	)

	return err
}
