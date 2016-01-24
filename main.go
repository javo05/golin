package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/gophergala2016/golin/login"
	"github.com/gophergala2016/golin/tokens"
	"time"
)

var secret = "ChAvO"

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

func (u User) Login(user, password string) (string, string) {
	if password == "qwerty" {
		return user, "chavo@smmx.in"
	}
	return "No", "NO"
}

func (t Token) GenerateToken(SignatureStr string) (string, string) {
	var claim Claim
	var tokenStr string

	claim.Id = SignatureStr
	claim.exp = time.Now().Add(time.Hour * 1).Unix()

	alg := jwt.GetSigningMethod("HS256")
	token := jwt.New(alg)
	token.Claims = structs.Map(claim)

	if tokenStr, err := token.SignedString([]byte(secret)); err == nil {
		fmt.Println(tokenStr) //TODO 1. aquí si lo genera
	} else {
		return "", fmt.Sprintf("Error signing token: %v", err)
	}
	//TODO 3. Aquí deberíamos implementar el storage
	return tokenStr, ""
}

func main() {
	r := gin.Default()
	publics := r.Group("api/v1/public")
	privates := r.Group("api/v1/private")
	//privates.Use(jwt.Auth(secret)) //TODO Change by verify method
	privates.GET("/users/:id", GetUser)
	publics.POST("/users", LoginUser)
	r.Run(":8080")
}

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
	//token := tokens.GenerateToken(SignatureStr)
	tokenStr, err := tokener.GenerateToken(SignatureStr)

	if err != "" {
		c.JSON(404, gin.H{"error": err})
	} else {
		c.JSON(201, gin.H{"token": tokenStr, "email": email})//TODO 2. Aquí no lo está imprimiendo
	}

}

func GetUser(c *gin.Context) {
	_ = c.Params.ByName("id")
	c.JSON(200, gin.H{"message": "chingon perron"})
}
