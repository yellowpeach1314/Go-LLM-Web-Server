package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
