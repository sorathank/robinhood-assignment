package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username string `bson:"_id" binding:"required"`
	Password string `bson:"password" binding:"required"`
	Email    string `bson:"email" binding:"required"`
}

type Login struct {
	Username string `bson:"username" binding:"required"`
	Password string `bson:"password" binding:"required"`
}

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetUserByUsername(username string) (User, error) {
	collection := ur.db.Collection("users")

	var user User
	err := collection.FindOne(context.TODO(), bson.D{
		{Key: "_id", Value: username}}).Decode(&user)

	return user, err
}

func (ur *UserRepository) GetUserByEmail(email string) (User, error) {
	collection := ur.db.Collection("users")

	var user User
	err := collection.FindOne(context.TODO(), bson.D{
		{Key: "email", Value: email}}).Decode(&user)

	return user, err
}

func (ur *UserRepository) InsertUser(user User) error {
	collection := ur.db.Collection("users")

	_, err := collection.InsertOne(context.TODO(), user)

	return err
}
