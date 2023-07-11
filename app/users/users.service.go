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

func ValidateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Body"})
			return
		}

		user, err := getUserByUsername(c, login.Username)
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

		c.JSON(http.StatusAccepted, gin.H{"Login": "Success"})
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
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create User": err.Error()})
			return
		}

		_, err := getUserByUsername(c, user.Username)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"Create User": "Duplicated Username"})
			return
		}

		_, err = getUserByEmail(c, user.Email)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"Create User": "Duplicated Email"})
			return
		}

		hash, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Fatal(err)
		}

		user.Password = hash
		insertErr := insertUser(c, user)
		if insertErr != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Create User": err.Error()})
			return
		}

		log.Println("Create User Success")

		c.JSON(http.StatusAccepted, gin.H{"Create User": "Success"})
	}
}
