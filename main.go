package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/configs"
	"github.com/sorathank/robinhood-assignment/app/interviews"
	"github.com/sorathank/robinhood-assignment/app/middleware"
	"github.com/sorathank/robinhood-assignment/app/users"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initConfig() configs.Configuration {
	var cf configs.Configuration
	v := viper.New()
	var configName = ""

	if os.Getenv("STAGE") == "" {
		configName = "config.dev"
	} else {
		configName = "config." + os.Getenv("STAGE")
	}
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath("./app/configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	err := v.Unmarshal(&cf)
	if err != nil {
		panic(err)
	}

	return cf
}

func connectToMongoDB(cf configs.Configuration) *mongo.Database {
	clientOptions := options.Client().ApplyURI(cf.MongoDb.Connection)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(cf.MongoDb.DatabaseName)
}

func setupRouter(cf configs.Configuration, db *mongo.Database) *gin.Engine {
	r := gin.Default()
	store, _ := redis.NewStore(10, "tcp", cf.Redis.Connection, "", []byte("secret"))
	sessionNames := []string{cf.Redis.SessionName.UserSession}
	r.Use(sessions.SessionsMany(sessionNames, store))

	sessionManager := middleware.NewSessionManager()
	authMiddleware := middleware.NewAuthMiddleware(sessionManager)

	users.UsersRoutes(r, db, cf)
	interviews.InterviewRoutes(r, db, cf, authMiddleware)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	cf := initConfig()
	db := connectToMongoDB(cf)
	r := setupRouter(cf, db)
	r.Run(":8080")
}
