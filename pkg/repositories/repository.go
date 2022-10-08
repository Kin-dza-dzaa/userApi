package repository

import (
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	SignUpUser(user *models.User) error
	VerifyUser(user *models.User) error
	GetUUid(refreshToken string) (string, error)
	GetVerifiedUser(user *models.User) ([2]string, error)
	UpdateRefreshToken(user *models.User, newToken string) error
} 

func NewRepository(pool *pgxpool.Pool) Repository {
	return NewUserRepository(pool)
}
