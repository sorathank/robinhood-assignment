package users

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sorathank/robinhood-assignment/app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	db         *mongo.Database
	repository *UserRepository
}

func NewUserController(db *mongo.Database) *UserController {
	repository := NewUserRepository(db)
	return &UserController{db: db, repository: repository}
}

func (uc *UserController) ValidateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Body"})
			return
		}

		user, err := uc.repository.GetUserByUsername(login.Username)
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

		uc.createUserSession(c, login.Username)
	}
}

func (uc *UserController) createUserSession(c *gin.Context, username string) {
	session := sessions.DefaultMany(c, "user_session")
	session.Set("username", username)
	session.Save()
	log.Println("User Session created")
	c.JSON(200, gin.H{"result": "Login Success"})
}

func (uc *UserController) CreateNewUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Create User": err.Error()})
			return
		}

		_, err := uc.repository.GetUserByUsername(user.Username)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"Create User": "Duplicated Username"})
			return
		}

		_, err = uc.repository.GetUserByEmail(user.Email)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"Create User": "Duplicated Email"})
			return
		}

		hash, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Fatal(err)
		}

		user.Password = hash
		insertErr := uc.repository.InsertUser(user)
		if insertErr != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Create User": err.Error()})
			return
		}

		log.Println("Create User Success")

		c.JSON(http.StatusAccepted, gin.H{"Create User": "Success"})
	}
}
