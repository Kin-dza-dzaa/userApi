package repository

import (
	"context"

	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	AddUser(context.Context, *models.User) error
	UpdateCredentials(context.Context, *models.User) error
	VerifyUser(context.Context, *models.User) error
	GetUUid(context.Context, *models.User) error
	GetVerifiedUser(context.Context, *models.User) (*models.User, error)
	UpdateRefreshToken(context.Context, *models.User) error
	IfUnverifiedUserExists(context.Context, *models.User, *bool) error
} 

func NewRepository(pool *pgxpool.Pool) Repository {
	return NewUserRepository(pool)
}

