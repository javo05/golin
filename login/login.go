package login

type Loginer interface {
    Login (user, password string) string
}

type SocialLoginer interface {
    Login (token string) string
}

