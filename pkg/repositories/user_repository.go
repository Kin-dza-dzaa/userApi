package repository

import (
	"context"
	"errors"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
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

type UserRepositry struct {
	logger *zerolog.Logger
	pool   *pgxpool.Pool
}

func (repository *UserRepositry) AddUser(user *models.User) error {
	if _, err := repository.pool.Exec(context.TODO(), queryCreateUser, user.UserId, user.UserName, user.Email, user.Password, user.RegistrationTime, user.VerificationCode, user.Verified); err != nil {
		return errors.New("user already exists")
	}
	return nil
}

func (repository *UserRepositry) UpdateCredentials(user *models.User) error {
	if _, err := repository.pool.Exec(context.TODO(), queryUpdateCreditnails, user.UserName, user.Password, user.VerificationCode, user.Email); err != nil {
		repository.logger.Error().Msg(err.Error())
		return errors.New("internal error")
	}
	return nil
}

func (repository *UserRepositry) VerifyUser(user *models.User) error {
	commandTag, err := repository.pool.Exec(context.TODO(), queryVerifyUser, user.RefreshToken, user.ExpirationTime, user.VerificationCode)
	if err != nil {
		repository.logger.Error().Msg(err.Error())
		return errors.New("internal error")
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("wrong verification code")
	}
	return nil
}

func (repository *UserRepositry) GetUUid(user *models.User) error {
	var userId string
	if err := repository.pool.QueryRow(context.TODO(), queryGetUUid, user.RefreshToken, time.Now().UTC()).Scan(&userId); err != nil {
		return errors.New("user doesn't exist")
	}
	UUid, err := uuid.Parse(userId)
	if err != nil {
		repository.logger.Error().Msg(err.Error())
		return errors.New("internal error")
	}
	user.UserId = UUid
	return nil
}

func (repository *UserRepositry) GetVerifiedUser(user *models.User) (*models.User, error) {
	var dbUser models.User
	if err := repository.pool.QueryRow(context.TODO(), queryGetVerifiedUser, user.Email).Scan(&dbUser.UserId, &dbUser.Password, &dbUser.RefreshToken, &dbUser.ExpirationTime); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("wrong email")
		}
		repository.logger.Error().Msg(err.Error())
		return nil, errors.New("internal error")
	}
	return &dbUser, nil
}

func (repository *UserRepositry) UpdateRefreshToken(user *models.User) error {
	_, err := repository.pool.Exec(context.TODO(), queryUpdateRefreshToken, user.RefreshToken, user.ExpirationTime, user.Email)
	if err != nil {
		repository.logger.Error().Msg(err.Error())
		return errors.New("internal error")
	}
	return nil
}

func (repository *UserRepositry) IfUnverifiedUserExists(user *models.User) (bool, error) {
	var result bool
	if err := repository.pool.QueryRow(context.TODO(), queryIfUnverifiedUserExists, user.Email).Scan(&result); err != nil {
		repository.logger.Error().Msg(err.Error())
		return false, errors.New("internal error")
	}
	return result, nil
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepositry {
	return &UserRepositry{
		pool: pool,
	}
}
