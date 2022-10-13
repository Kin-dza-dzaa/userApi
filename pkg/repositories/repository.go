package repository

import (
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Repository interface {
	AddUser(user *models.User) error
	UpdateCredentials(user *models.User) error
	VerifyUser(user *models.User) error
	GetUUid(user *models.User) error
	GetVerifiedUser(user *models.User) (*models.User, error)
	UpdateRefreshToken(user *models.User) error
	IfUnverifiedUserExists(user *models.User) (bool, error)
} 

func NewRepository(pool *pgxpool.Pool, logger *zerolog.Logger) Repository {
	return NewUserRepository(pool)
}
