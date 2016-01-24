package main

import (
    "github.com/gin-gonic/gin"
    "github.com/goris/golin/login"
    "github.com/goris/golin/tokens"
    "github.com/gin-gonic/contrib/jwt"
)

var secret = "ChAvO"

type User struct {
    AccountId string `json:"account_id,omitempty"`
    Username string `json:"username,omitempty"`
    Email string `json:"email,omitempty"`
    Password string `json:"password,omitempty"`
}

func (u User) Login(user, password string) (string, string) {
    if password == "qwerty" {
        return user, "chavo@smmx.in"
    }
    return "No", "NO"
} 

func main() {
    r := gin.Default()
    publics := r.Group("api/v1/public")
    privates := r.Group("api/v1/private")
    privates.Use(jwt.Auth(secret))
    {
        privates.GET("/users/:id", GetUser)
        publics.POST("/users", LoginUser)
    }
    r.Run(":8080")
}

func LoginUser(c *gin.Context) {
    var user User
    var loginer login.Loginer

    c.BindJSON(&user)

    loginer = user

    SignatureStr, email := loginer.Login(user.Username, user.Password)
    token := tokens.GenerateToken(SignatureStr)
    c.JSON(201, gin.H{"token": token, "email": email})

}

func GetUser(c *gin.Context) {
    _ = c.Params.ByName("id")
    c.JSON(200, gin.H{"message": "chingon perron"})
}
