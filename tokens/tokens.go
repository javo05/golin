package tokens

type Tokener interface {
	GenerateToken(SignaturStr string) string
}

var secret = "ChAvO"
