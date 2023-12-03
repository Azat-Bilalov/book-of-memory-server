package ds

import (
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	jwt.StandardClaims
	User_id string `json:"user_id"`
	Role    string `json:"role"`
}
