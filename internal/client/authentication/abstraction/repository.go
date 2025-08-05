package abstraction

type authRepository struct {
	token string
	login string
}

func NewRepository() *authRepository {
	return &authRepository{}
}

func (r *authRepository) GetUser() string {
	return r.login
}

func (r *authRepository) GetToken() string {
	return r.token
}

func (r *authRepository) Authenticate(login, token string) {
	r.login = login
	r.token = token
}
