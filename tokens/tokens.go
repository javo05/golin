package tokens

import (
    jwt_lib "github.com/dgrijalva/jwt-go"
    "time"
)

var secret = "ChAvO"

func GenerateToken(SignatureStr string) string {
    var tokenStr string

    token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
    token.Claims["ID"] = SignatureStr
    token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
    tokenStr, _ = token.SignedString([]byte(secret))

    return tokenStr
}
