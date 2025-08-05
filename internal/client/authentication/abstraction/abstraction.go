package abstraction

type Repository interface {
	GetUser() string
	GetToken() string
	Authenticate(login string, token string)
}
