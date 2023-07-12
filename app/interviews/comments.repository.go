package interviews

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) *CommentRepository {
	return &CommentRepository{
		collection: db.Collection("comments"),
	}
}
func (r *CommentRepository) FindByInterviewID(interviewId string) ([]Comment, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_time", Value: -1}})
	cursor, err := r.collection.Find(context.Background(), bson.M{"interview_id": interviewId}, findOptions)

	var comments []Comment
	if err = cursor.All(context.Background(), &comments); err != nil {
		log.Println(err)
	}

	return comments, err
}
func (r *CommentRepository) Insert(comment Comment, user string) error {
	insertResult, err := r.collection.InsertOne(context.Background(), bson.M{
		"User":        user,
		"InterviewId": comment.InterviewId,
		"Content":     comment.Content,
		"CreatedTime": time.Now(),
	})
	fmt.Println("Inserted new comment: ", insertResult.InsertedID)
	return err
}
