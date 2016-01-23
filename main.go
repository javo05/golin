package main

import (
    "github.com/gin-gonic/gin"
    "github.com/goris/golin/login"
    "fmt"
)

type User struct {
    AccountId string `json:"account_id,omitempty"`
    Username string `json:"username,omitempty"`
    Email string `json:"email,omitempty"`
    Password string `json:"password,omitempty"`
}

func (u User) Login(user, password string) string {
    if user == "javo" && password == "qwerty" {
        return "12345"
    }
    return "No"
} 

func main() {
    r := gin.Default()
    v1 := r.Group("api/v1")
    {
        v1.GET("/users/:id", GetUser)
        v1.POST("/users", LoginUser)
    }
    r.Run(":8080")
}

func LoginUser(c *gin.Context) {
    var user User
    var loginer login.Loginer

    c.BindJSON(&user)


    loginer = user

    fmt.Println(user.Username, user.Password)
    token := loginer.Login(user.Username, user.Password)
    c.JSON(201, token)

}

func GetUser(c *gin.Context) {
}
