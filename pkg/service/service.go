package service

import (
	"context"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/repositories"
)

type Service interface {
	SignUpUser(ctx context.Context, user *models.User) error
	SignInUser(ctx context.Context, user *models.User) error
	VerifyUser(ctx context.Context, user *models.User) error
	GetAccessToken(ctx context.Context, user *models.User) error
}

func NewService(repository repository.Repository, config *config.Config) Service {
	return NewUserService(repository, config)
}