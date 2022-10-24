package repository

import (
	"context"

	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	AddUser(ctx context.Context, user *models.User) error
	UpdateCredentials(ctx context.Context, user *models.User) error
	VerifyUser(ctx context.Context, user *models.User) error
	GetUUid(ctx context.Context, user *models.User) error
	GetVerifiedUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateRefreshToken(ctx context.Context, user *models.User) error
	IfUnverifiedUserExists(ctx context.Context, user *models.User) (bool, error)
} 

func NewRepository(pool *pgxpool.Pool) Repository {
	return NewUserRepository(pool)
}
