package repository

import (
	"context"
	"errors"
	"time"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepositry struct {
	pool *pgxpool.Pool
}

func (a *UserRepositry) SignUpUser(user *models.User) error {
	sqlCommand := `INSERT INTO users(id, user_name, email, password, registration_date, verification_code, verified) 
						VALUES($1, $2, $3, $4, $5, $6, $7);`
	if _, err := a.pool.Exec(context.TODO(), sqlCommand, user.UserId, user.UserName, user.Email, user.Password, user.RegistrationTime, user.VerificationCode, user.Verified); err != nil {
		return errors.New("user already exists")
	}
	return nil
}

func (repository *UserRepositry) VerifyUser(user *models.User) error {
	sqlCommand := `UPDATE users SET 
					verified = true, refresh_token = $1, expiration_time = $2 WHERE verification_code = $3;`
	commandTag, _ := repository.pool.Exec(context.TODO(), sqlCommand, user.RefreshToken, user.ExpirationTime, user.VerificationCode)
	if commandTag.RowsAffected() == 0 {
		return errors.New("wrong verification code")
	}
	return nil
}

func (repository *UserRepositry) GetUUid(refreshToken string) (string, error) {
	sqlCommand := `SELECT id FROM users WHERE
					refresh_token = $1 AND 
					verified = true AND
					expiration_time > $2;
					`
	var userId string
	if err := repository.pool.QueryRow(context.TODO(), sqlCommand, refreshToken, time.Now().UTC()).Scan(&userId); err != nil {
		return "", errors.New("user doesn't exist")
	}
	return userId, nil
}

func (repository *UserRepositry) GetVerifiedUser(user *models.User) ([2]string, error) {
	sqlCommand := `SELECT id, password FROM users WHERE
					email = $1 AND
					verified = true;
	`
	var UUidString string
	var password string
	if err := repository.pool.QueryRow(context.TODO(), sqlCommand, user.Email).Scan(&UUidString, &password); err != nil {
		return [2]string{}, errors.New("internal error")
	}
	return [2]string{UUidString, password}, nil
}

func (repostiory *UserRepositry) UpdateRefreshToken(user *models.User, newToken string) error {
	sqlCommand := `UPDATE users SET refresh_token = $1, expiration_time = $2 WHERE refresh_token = $3;`
	_, err := repostiory.pool.Exec(context.TODO(), sqlCommand, newToken, user.ExpirationTime, user.RefreshToken)
	if err != nil {
		return errors.New("internal error")
	}
	return nil
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepositry {
	return &UserRepositry{
		pool: pool,
	}
}