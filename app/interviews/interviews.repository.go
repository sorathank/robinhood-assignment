package interviews

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type InterviewRepository struct {
	collection *mongo.Collection
}

func NewInterviewRepository(db *mongo.Database) *InterviewRepository {
	return &InterviewRepository{
		collection: db.Collection("interviews"),
	}
}

func (r *InterviewRepository) FindOneByID(id string) (Interview, error) {
	var interview Interview
	objectId, err := primitive.ObjectIDFromHex(id)
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&interview)
	return interview, err
}

func (r *InterviewRepository) FindWithPagination(pageSize int64, pageNumber int64) ([]Interview, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_time", Value: -1}})
	findOptions.SetSkip((pageNumber - 1) * pageSize)
	findOptions.SetLimit(pageSize)

	cursor, err := r.collection.Find(context.Background(), bson.M{"status": bson.M{"$ne": Archived}}, findOptions)

	var interviews []Interview
	if err = cursor.All(context.Background(), &interviews); err != nil {
		log.Println(err)
	}

	return interviews, err
}

func (r *InterviewRepository) Insert(interview Interview, creator string) error {
	log.Println(time.Now())
	insertResult, err := r.collection.InsertOne(context.Background(), bson.M{
		"description":  interview.Description,
		"user":         creator,
		"status":       Todo,
		"created_time": time.Now(),
	})
	fmt.Println("Inserted new Interview: ", insertResult.InsertedID)
	return err
}

func (r *InterviewRepository) UpdateStatus(status Status, interviewId string) error {
	objectId, err := primitive.ObjectIDFromHex(interviewId)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectId},
		bson.D{
			{"$set", bson.D{{"status", status}}},
		},
	)
	return err
}
