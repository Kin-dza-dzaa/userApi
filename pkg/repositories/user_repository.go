package repository

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Kin-dza-dzaa/userApi/internal/apierror"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
	ErrUserAlredyExists = errors.New("user already exists")
	ErrUserDoesntExists = errors.New("refresh token expired or user doesn't exists")
	ErrWrongEmail = errors.New("wrong email")
)

type UserRepositry struct {
	pool dbconn
}

func (repository *UserRepositry) AddUser(ctx context.Context, user *models.User) error {
	if _, err := repository.pool.Exec(ctx, queryCreateUser, user.UserId, user.UserName, user.Email, user.Password, user.RegistrationTime, user.VerificationCode, user.Verified); err != nil {
		return apierror.NewErrorStruct(ErrUserAlredyExists.Error(), "error", http.StatusBadRequest)
	}
	return nil
}

func (repository *UserRepositry) UpdateCredentials(ctx context.Context, user *models.User) error {
	if _, err := repository.pool.Exec(ctx, queryUpdateCreditnails, user.UserName, user.Password, user.VerificationCode, user.Email); err != nil {
		return apierror.NewErrorStruct(ErrUserAlredyExists.Error(), "error", http.StatusBadRequest)
	}
	return nil
}

func (repository *UserRepositry) VerifyUser(ctx context.Context, user *models.User) error {
	commandTag, err := repository.pool.Exec(ctx, queryVerifyUser, user.RefreshToken, user.ExpirationTime, user.VerificationCode)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return apierror.NewErrorStruct(ErrWrongVerificationCode.Error(), "error", http.StatusBadRequest)
	}
	return nil
}

func (repository *UserRepositry) GetUUid(ctx context.Context, user *models.User) error {
	var userId string
	if err := repository.pool.QueryRow(ctx, queryGetUUid, user.RefreshToken, time.Now().UTC()).Scan(&userId); err != nil {
		if err == pgx.ErrNoRows {
			return apierror.NewErrorStruct(ErrUserDoesntExists.Error(), "error", http.StatusBadRequest)
		}
		return err
	}
	UUid, err := uuid.Parse(userId)
	if err != nil {
		return apierror.NewErrorStruct(ErrUserDoesntExists.Error(), "error", http.StatusBadRequest)
	}
	user.UserId = UUid
	return nil
}

func (repository *UserRepositry) GetVerifiedUser(ctx context.Context, user *models.User) (*models.User, error) {
	var dbUser models.User
	if err := repository.pool.QueryRow(ctx, queryGetVerifiedUser, user.Email).Scan(&dbUser.UserId, &dbUser.Password, &dbUser.RefreshToken, &dbUser.ExpirationTime); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apierror.NewErrorStruct(ErrWrongEmail.Error(), "error", http.StatusBadRequest)
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

func (repository *UserRepositry) IfUnverifiedUserExists(ctx context.Context, user *models.User, result *bool) (error) {
	if err := repository.pool.QueryRow(ctx, queryIfUnverifiedUserExists, user.Email).Scan(result); err != nil {
		return err
	}
	return nil
}

func NewUserRepository(pool dbconn) *UserRepositry {
	return &UserRepositry{
		pool: pool,
	}
}

// for tests
type dbconn interface {
	Close()
    Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
    QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
    SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
    Begin(ctx context.Context) (pgx.Tx, error)
    BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
    BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
    BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
}