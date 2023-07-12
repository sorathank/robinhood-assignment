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
	"github.com/sorathank/robinhood-assignment/app/users"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db = make(map[string]string)

func mongoPool(cf configs.Configuration) gin.HandlerFunc {
	clientOptions := options.Client().ApplyURI(cf.MongoDb.Connection)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return func(c *gin.Context) {
		c.Set(cf.MongoDb.SessionName, client.Database(cf.MongoDb.DatabaseName))
		c.Next()
	}
}

func setupRouter() *gin.Engine {
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

	r := gin.Default()
	store, _ := redis.NewStore(10, "tcp", cf.Redis.Connection, "", []byte("secret"))
	sessionNames := []string{cf.Redis.SessionName.UserSession}
	r.Use(sessions.SessionsMany(sessionNames, store))
	r.Use(mongoPool(cf))
	// Disable Console Color
	// gin.DisableConsoleColor()
	// r := gin.Default()

	users.UsersRoutes(r)
	interviews.InterviewRoutes(r, cf)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8000")
}
