package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserId    		 	uuid.UUID
	UserName  		 	string
	Email     		 	string   
	Password  		 	string    	
	RegistrationTime 	time.Time   
	RefreshToken 		string		
	ExpirationTime 		time.Time   
	VerificationCode	string		
	Verified		 	bool 		
	CsrfToken			string 		
	Jwt					string		
}
