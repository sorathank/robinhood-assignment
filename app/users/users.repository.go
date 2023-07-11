package users

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getDB(c *gin.Context) *mongo.Database {
	return c.MustGet("mongo_session").(*mongo.Database)
}

func getUser(c *gin.Context, username string) (Login, error) {
	db := getDB(c)
	collection := db.Collection("Users")

	var login Login
	err := collection.FindOne(context.TODO(), bson.D{
		{Key: "username", Value: username}}).Decode(&login)

	return login, err
}

func insertUser(c *gin.Context, hashedUser Login) error {
	db := getDB(c)
	collection := db.Collection("Users")

	insertResult, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "username", Value: hashedUser.Username},
		{Key: "password", Value: hashedUser.Password}})
	fmt.Println("Inserted new user: ", insertResult.InsertedID)

	return err
}
