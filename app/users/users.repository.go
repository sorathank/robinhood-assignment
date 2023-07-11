package users

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	UserId   primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Email    string             `bson:"email"`
}

type Login struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func getDB(c *gin.Context) *mongo.Database {
	return c.MustGet("mongo_session").(*mongo.Database)
}

func getUser(c *gin.Context, username string) (Login, error) {
	db := getDB(c)
	collection := db.Collection("users")

	var login Login
	err := collection.FindOne(context.TODO(), bson.D{
		{Key: "username", Value: username}}).Decode(&login)

	return login, err
}

func insertUser(c *gin.Context, hashedUser Login) error {
	db := getDB(c)
	collection := db.Collection("users")

	insertResult, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "username", Value: hashedUser.Username},
		{Key: "password", Value: hashedUser.Password}})
	fmt.Println("Inserted new user: ", insertResult.InsertedID)

	return err
}
