package users

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/utils"
)

// type UserController struct {
// 	CF configs.Configuration
// }

type Login struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func ValidateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			log.Println(err)
			log.Println("TEST")

			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Body"})
			return
		}

		log.Println("TEST")

		user, err := getUser(c, login.Username)
		var errorMessage interface{} = "Username or Password is incorrect"
		if err != nil {
			//User not found
			c.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
			return
		}

		isPasswordCorrect := utils.CheckPasswordWithHash(login.Password, user.Password)
		if !isPasswordCorrect {
			//Invalid password
			c.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
			return
		}

		log.Println("Login Success")
		createUserSession(c, login.Username)
	}
}

func createUserSession(c *gin.Context, username string) {
	session := sessions.DefaultMany(c, "user_session")
	session.Set("username", username)
	session.Save()
	log.Println("User Session created")
	c.JSON(200, gin.H{"result": "Login Success"})
}

func CreateNewUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user Login
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := getUser(c, user.Username)
		if err == nil {
			var errorMessage interface{} = "Duplicated Username"
			c.JSON(http.StatusConflict, gin.H{"error": errorMessage})
			return
		}

		log.Println("TEST")
		log.Println(user.Username)
		log.Println(user.Password)

		hash, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Fatal(err)
		}

		user.Password = hash
		insertErr := insertUser(c, user)
		if insertErr != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("Create User Success")
	}
}
