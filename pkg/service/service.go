package service

import (
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	SignUpUser(user *models.User) error
	SignInUser(user *models.User) (string, error)
	VerifyUser(user *models.User) (string, error)
	GetAccessToken(refreshToken string) (string, error)
}

func NewService(repository repository.Repository, config *config.Config, validator *validator.Validate) Service {
	return NewUserService(repository, config, validator)
}