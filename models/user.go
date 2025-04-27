package models

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	FirstName    string `json:"fistName"`
	LastName     string `json:"lastName"`
}

type RegisterUserReqBody struct {
	FirstName string `json:"fistName" binding:"required"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=7"`
}
