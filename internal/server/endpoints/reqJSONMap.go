package endpoints

type LoginRequest struct {
	Name     string
	Password string
}

type RegisterRequest struct {
	Name     string
	Password string
}

type LoginResponse struct {
	Status int64
	Token  string
}

type RegisterResponse struct {
}

type SaveRequest struct {
	Name string
	Data []byte
}

type SaveResponse struct {
	Size uint32
}
