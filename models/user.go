package models

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	PasswordHash string `json:"password;omitempty"`
}