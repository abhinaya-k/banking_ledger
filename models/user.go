package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	FirstName    string `json:"fistName"`
	LastName     string `json:"lastName"`
	Role         string `json:"role"`
}

type RegisterUserReqBody struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=7"`
	Role      string `json:"role" binding:"required,oneof=admin user"`
}

type LoginRequestBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseBody struct {
	Token string `json:"token" binding:"required"`
}

type CustomClaims struct {
	Role string `json:"role"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}
