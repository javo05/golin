package tokens

/*
/   Simple interface for using creating a Token according to
/   different needs
*/

type Tokener interface {
	GenerateToken(SignaturStr string) (string, error)
}

var secret = "ChAvO"
