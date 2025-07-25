package dto

type User struct {
	Id       int64
	Name     string
	Password string
}

type UserCtx string

const (
	Id UserCtx = "userId"
)
