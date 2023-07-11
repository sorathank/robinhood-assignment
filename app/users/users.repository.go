package users

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username string `bson:"_id"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

type Login struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func getDB(c *gin.Context) *mongo.Database {
	return c.MustGet("mongo_session").(*mongo.Database)
}

func getUserByUsername(c *gin.Context, username string) (User, error) {
	db := getDB(c)
	collection := db.Collection("users")

	var user User
	err := collection.FindOne(context.TODO(), bson.D{
		{Key: "_id", Value: username}}).Decode(&user)

	return user, err
}

func getUserByEmail(c *gin.Context, email string) (User, error) {
	db := getDB(c)
	collection := db.Collection("users")

	var user User
	err := collection.FindOne(context.TODO(), bson.D{
		{Key: "email", Value: email}}).Decode(&user)

	return user, err
}

func insertUser(c *gin.Context, user User) error {
	db := getDB(c)
	collection := db.Collection("users")

	insertResult, err := collection.InsertOne(context.TODO(), user)
	fmt.Println("Inserted new user: ", insertResult.InsertedID)

	return err
}
