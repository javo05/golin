package main

import (
	"encoding/json"
	"fmt"
    "github.com/boltdb/bolt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
    "github.com/gophergala2016/golin/boltdb"
	"github.com/gophergala2016/golin/login"
	"github.com/gophergala2016/golin/tokens"
    "log"
	"regexp"
	"strings"
	"time"
)

// This secret should be included via a configuration file
// The configuration file should also include to which DB it should
// connnect to retrieve and ensure the login of the users.
// ANother things that the configuration file should include are
// the tables where the users information are and also which 
// columns to look into.
var secret = "ChAvO"

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

type Token struct {
	Token string
}

type Claim struct {
	Id  string
	exp int64
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

// JWT way to create and generate tokens
func (t Token) GenerateToken(SignatureStr string) (string, error) {
	var claim Claim

	claim.Id = SignatureStr
	claim.exp = time.Now().Add(time.Hour * 1).Unix()

	alg := jwt.GetSigningMethod("HS256")
	token := jwt.New(alg)
	token.Claims = structs.Map(claim)

	if tokenStr, err := token.SignedString([]byte(secret)); err == nil {
		return tokenStr, nil
	} else {
		return "", err
	}
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

func ValidateToken(encriptedToken string) (string, error) {
	tokData := regexp.MustCompile(`\s*$`).ReplaceAll([]byte(encriptedToken), []byte{})

	currentToken, err := jwt.Parse(string(tokData), func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// Print an error if we can't parse for some reason
	if err != nil {
		fmt.Println("Couldn't parse token: ", err)
		return "", fmt.Errorf("Couldn't parse token: %v", err)
	}

	fmt.Println("currentToken ", currentToken)
	// Is token invalid?
	if !currentToken.Valid {
		return "", fmt.Errorf("Token is invalid")
	}

	if err != nil {
		//Validate if token has expired or is not parseable
		//Validate in the blacklist
		return "", err
	}

	// Print the token details
	_, err = json.MarshalIndent(currentToken.Claims, "", "    ")

	return string(tokData), err
}

// These are the endpoints required to do a login and verifying that tokens are
// alive 
func LoginUser(c *gin.Context) {
	var user User
	var loginer login.Loginer
	var tokenStr string

	var token Token
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

	newToken, err := ValidateToken(currentToken)

	if err != nil {
		c.JSON(404, gin.H{"error": err})
	} else {
		c.JSON(200, gin.H{"message": "chingon perron", "token": newToken})
	}
}
