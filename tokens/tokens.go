package tokens

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"regexp"
	"time"
)

// This secret should be included via a configuration file
// The configuration file should also include to which DB it should
// connnect to retrieve and ensure the login of the users.
// ANother things that the configuration file should include are
// the tables where the users information are and also which
// columns to look into.
var secret = "ChAvO"

/*
/   Simple interface for using creating a Token according to
/   different needs
*/

type Tokener interface {
	GenerateToken(SignaturStr string) (string, error)
}

type Token struct {
	Token string
}

type Claim struct {
	Id  string
	exp int64
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

/*
* Validate whether the token is a valid JWT format
 */
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
