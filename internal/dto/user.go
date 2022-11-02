package dto

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/Kin-dza-dzaa/userApi/internal/apierror"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSignInDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto UserSignInDto) IntoUser() (*models.User, error) {
	email, err := mail.ParseAddress(dto.Email)
	if err != nil {
		return nil, apierror.NewErrorStruct(ErrInvalidCredentials.Error(), "error", http.StatusBadRequest)
	}
	if len(dto.Password) < 8 {
		return nil, apierror.NewErrorStruct(ErrInvalidCredentials.Error(), "error", http.StatusBadRequest)
	}
	var User models.User
	User.Email = email.Address
	User.Password = dto.Password
	return &User, nil
}

type UserSignUpDto struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto UserSignUpDto) IntoUser() (*models.User, error) {
	email, err := mail.ParseAddress(dto.Email)
	if err != nil {
		return nil, apierror.NewErrorStruct(ErrInvalidCredentials.Error(), "error", http.StatusBadRequest)
	}

	if len(dto.Password) < 8 {
		return nil, apierror.NewErrorStruct(ErrInvalidCredentials.Error(), "error", http.StatusBadRequest)
	}
	if len(dto.UserName) < 8 {
		return nil, apierror.NewErrorStruct(ErrInvalidCredentials.Error(), "error", http.StatusBadRequest)
	}
	var User models.User
	User.Email = email.Address
	User.UserName = dto.UserName
	User.Password = dto.Password
	return &User, nil
}
