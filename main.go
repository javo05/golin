package main

import (
	"github.com/boltdb/bolt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/goris/golin/boltdb"
	"github.com/goris/golin/login"
	"github.com/goris/golin/tokens"
	"log"
	"strings"
)

// BoltDB was used because of it's speed and architecture goes
// very well accordingly into what we want to achieve, which is
// a secure and fast service to consume tokens.
var db *bolt.DB

type User struct {
	AccountId string `json:"account_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

// This is where config file should be used to read and compare users
// since this is the MVP of this microservice, this works for achieving
// what we want.
func (u User) Login(user, password string) (string, string) {
	if password == "qwerty" {
		return user, "chavo@segundamano.mx"
	}
	return "No", "NO"
}

// Use config file to decide which DB you'll be using for storing the tokens
func main() {
	var err error
	r := gin.Default()
	db, err = boltdb.OpenBoltDB("tokens")
	if err != nil {
		log.Fatal(err)
	}
	publics := r.Group("api/v1/public")
	privates := r.Group("api/v1/private")
	privates.GET("/users/:id", GetUser)
	publics.POST("/users", LoginUser)
	r.Run(":8080")
}

// These are the endpoints required to do a login and verifying that tokens are
// alive
func LoginUser(c *gin.Context) {
	var user User
	var loginer login.Loginer
	var tokenStr string

	var token tokens.Token
	var tokener tokens.Tokener

	c.BindJSON(&user)

	loginer = user
	tokener = token

	SignatureStr, email := loginer.Login(user.Username, user.Password)
	tokenStr, err := tokener.GenerateToken(SignatureStr)

	if err != nil {
		c.JSON(404, gin.H{"error generating token": err})
	} else {
		data := structs.Map(user)
		err = boltdb.UpdateBucket(db, tokenStr, data)
		if err != nil {
			c.JSON(404, gin.H{"error updating bucket": err})
		} else {
			c.JSON(201, gin.H{"token": tokenStr, "email": email})
		}
	}
}

// Here we use the function to validate if a token exists and if the user should
// alllowed to enter and go through
func GetUser(c *gin.Context) {
	_ = c.Params.ByName("id")
	var currentToken string

	authorization_array := c.Request.Header["Authorization"]
	if len(authorization_array) > 0 {
		currentToken = strings.Fields(strings.TrimSpace(authorization_array[0]))[1]
	}

	//Before this, token should not be in a black list
	newToken, err := tokens.ValidateToken(currentToken)

	if err != nil {
		c.JSON(404, gin.H{"error": err})
	} else {
		c.JSON(200, gin.H{"message": "chingon perron", "token": newToken})
	}
}
