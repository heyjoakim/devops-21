package models

type UserCreateRequest struct {
	Username  string
	Email     string
	Password  string
	Password2 string
}
