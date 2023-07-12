package utils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordWithHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "user_session")
		response := session.Get("username")

		c.JSON(200, gin.H{"username": response})
		fmt.Println(response)
	}
}

func clearUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "user_session")
		session.Clear()
		session.Save()
		c.JSON(200, gin.H{"success": 1})
	}
}

func GetDB(c *gin.Context) *mongo.Database {
	return c.MustGet("mongo_session").(*mongo.Database)
}

func StringToPositiveInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	if i <= int64(0) {
		return 0, errors.New("number is not positive")
	}
	return i, nil
}
