package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserId    		 	uuid.UUID	 `json:"-"`
	UserName  		 	string    	 `json:"user_name,omitempty" validate:"omitempty,userval"`
	Email     		 	string    	 `json:"email,omitempty" validate:"omitempty,email"`
	Password  		 	string    	 `json:"password" validate:"required,passval"`
	RegistrationTime 	time.Time    `json:"-"`
	RefreshToken 		string		 `json:"-"`
	ExpirationTime 		time.Time    `json:"-"`
	VerificationCode	string		 `json:"-"`
	Verified		 	bool 		 `json:"-"`
}
