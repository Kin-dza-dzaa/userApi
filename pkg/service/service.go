package service

import (
	"context"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/repositories"
)

type Service interface {
	SignUpUser(context.Context, *models.User) error
	SignInUser(context.Context, *models.User) error
	VerifyUser(context.Context, *models.User) error
	GetAccessToken(context.Context, *models.User) error
}

func NewService(repository repository.Repository, config *config.Config) Service {
	return NewUserService(repository, config)
}