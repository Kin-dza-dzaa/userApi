package service

import (
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Service interface {
	SignUpUser(user *models.User) error
	SignInUser(user *models.User) error
	VerifyUser(user *models.User) error
	GetAccessToken(user *models.User) error
}

func NewService(repository repository.Repository, config *config.Config, validator *validator.Validate, logger *zerolog.Logger) Service {
	return NewUserService(repository, config, validator, logger)
}