package tokens

type Tokener interface {
	GenerateToken(SignaturStr string) (string, error)
}

var secret = "ChAvO"
