package models

import "github.com/golang-jwt/jwt"

type MyJwtClaims struct {
	UserId 					string 		`json:"user_id"`
	XCSRFToken				string 		`json:"x_csrf_token"`
	jwt.StandardClaims
}