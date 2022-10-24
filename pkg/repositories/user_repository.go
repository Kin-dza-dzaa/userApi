package repository

import (
	"context"
	"errors"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

const (
	queryCreateUser             = "INSERT INTO users(id, user_name, email, password, registration_date, verification_code, verified) VALUES($1, $2, $3, $4, $5, $6, $7);"
	queryUpdateCreditnails      = "UPDATE users SET user_name = $1, password = $2, verification_code = $3 WHERE email = $4;"
	queryVerifyUser 	        = "UPDATE users SET verified = true, refresh_token = $1, expiration_time = $2, verification_code = '' WHERE verification_code = $3 AND verified = false;" 
	queryGetUUid		        = "SELECT id FROM users WHERE refresh_token = $1 AND verified = true AND expiration_time > $2;"
	queryGetVerifiedUser        = "SELECT id, password, refresh_token, expiration_time FROM users WHERE email = $1 AND verified = true;"
	queryUpdateRefreshToken     = "UPDATE users SET refresh_token = $1, expiration_time = $2 WHERE email = $3;"
	queryIfUnverifiedUserExists = "SELECT EXISTS(SELECT * FROM USERS WHERE email = $1 AND verified = false);"
)

var (
	ErrWrongVerificationCode = errors.New("wrong verification code")
	ErrUserDoesntExists = errors.New("user doesn't exist")
	ErrWrongEmail = errors.New("wrong email")
)

type UserRepositry struct {
	pool   *pgxpool.Pool
}

func (repository *UserRepositry) AddUser(ctx context.Context, user *models.User) error {
	if _, err := repository.pool.Exec(ctx, queryCreateUser, user.UserId, user.UserName, user.Email, user.Password, user.RegistrationTime, user.VerificationCode, user.Verified); err != nil {
		return err
	}
	return nil
}

func (repository *UserRepositry) UpdateCredentials(ctx context.Context, user *models.User) error {
	if _, err := repository.pool.Exec(ctx, queryUpdateCreditnails, user.UserName, user.Password, user.VerificationCode, user.Email); err != nil {
		return err
	}
	return nil
}

func (repository *UserRepositry) VerifyUser(ctx context.Context, user *models.User) error {
	commandTag, err := repository.pool.Exec(ctx, queryVerifyUser, user.RefreshToken, user.ExpirationTime, user.VerificationCode)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrWrongVerificationCode
	}
	return nil
}

func (repository *UserRepositry) GetUUid(ctx context.Context, user *models.User) error {
	var userId string
	if err := repository.pool.QueryRow(ctx, queryGetUUid, user.RefreshToken, time.Now().UTC()).Scan(&userId); err != nil {
		return ErrUserDoesntExists
	}
	UUid, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	user.UserId = UUid
	return nil
}

func (repository *UserRepositry) GetVerifiedUser(ctx context.Context, user *models.User) (*models.User, error) {
	var dbUser models.User
	if err := repository.pool.QueryRow(ctx, queryGetVerifiedUser, user.Email).Scan(&dbUser.UserId, &dbUser.Password, &dbUser.RefreshToken, &dbUser.ExpirationTime); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrWrongEmail
		}
		return nil, err
	}
	return &dbUser, nil
}

func (repository *UserRepositry) UpdateRefreshToken(ctx context.Context, user *models.User) error {
	_, err := repository.pool.Exec(ctx, queryUpdateRefreshToken, user.RefreshToken, user.ExpirationTime, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (repository *UserRepositry) IfUnverifiedUserExists(ctx context.Context, user *models.User) (bool, error) {
	var result bool
	if err := repository.pool.QueryRow(ctx, queryIfUnverifiedUserExists, user.Email).Scan(&result); err != nil {
		return false, err
	}
	return result, nil
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepositry {
	return &UserRepositry{
		pool: pool,
	}
}
