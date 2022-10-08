package models

import "github.com/golang-jwt/jwt"

type MyJwtClaims struct {
	UserId 					string 		`json:"user_id"`
	jwt.StandardClaims
}