package main

import (
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
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

func (u User) Login(user, password string) (string, string) {
	if password == "qwerty" {
		return user, "chavo@smmx.in"
	}
	return "No", "NO"
}

func (t Token) GenerateToken(SignatureStr string) string {
	var tokenStr string

	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	token.Claims["ID"] = SignatureStr
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	tokenStr, _ = token.SignedString([]byte(secret))

	return tokenStr
}

func main() {
	r := gin.Default()
	publics := r.Group("api/v1/public")
	privates := r.Group("api/v1/private")
	privates.Use(jwt.Auth(secret))
	privates.GET("/users/:id", GetUser)
	publics.POST("/users", LoginUser)
	r.Run(":8080")
}

func LoginUser(c *gin.Context) {
	var user User
	var loginer login.Loginer

	var token Token
	var tokener tokens.Tokener

	c.BindJSON(&user)

	loginer = user
	tokener = token

	SignatureStr, email := loginer.Login(user.Username, user.Password)
	//token := tokens.GenerateToken(SignatureStr)
	c.JSON(201, gin.H{"token": tokener.GenerateToken(SignatureStr), "email": email})

}

func GetUser(c *gin.Context) {
	_ = c.Params.ByName("id")
	c.JSON(200, gin.H{"message": "chingon perron"})
}
