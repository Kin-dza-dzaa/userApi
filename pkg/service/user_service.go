package service

import (
	"errors"
	"time"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	repository "github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/dchest/uniuri"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repository.Repository
	config     *config.Config
	validator  *validator.Validate
}

func (service *UserService) SignUpUser(user *models.User) error {
	if err := service.validator.Struct(user); err != nil {
		return errors.New("invalid creditnails")
	}
	if err := service.hashPassword(user); err != nil {
		return err
	}
	user.UserId = uuid.New()
	user.RegistrationTime = time.Now().UTC()
	user.VerificationCode = uniuri.New()
	user.Verified = false
	if err := service.repository.SignUpUser(user); err != nil {
		return err
	}
	// send email
	return nil
}

func (service *UserService) SignInUser(user *models.User) (string, error) {
	data, err := service.repository.GetVerifiedUser(user)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(data[1]), []byte(user.Password)); err != nil {
		return "", errors.New("invalid password")
	}
	user.RefreshToken = uniuri.NewLen(512)
	user.ExpirationTime = time.Now().UTC().AddDate(0, 6, 0)
	if err := service.repository.UpdateRefreshToken(user, user.RefreshToken); err != nil {
		return "", err
	}
	return service.generateToken(data[0])
}

func (service *UserService) VerifyUser(user *models.User) (string, error) {
	user.RefreshToken = uniuri.NewLen(512)
	user.ExpirationTime = time.Now().UTC().AddDate(0, 6, 0)
	if err := service.repository.VerifyUser(user); err != nil {
		return "", err
	}
	UUidString, err := service.repository.GetUUid(user.RefreshToken)
	if err != nil {
		return "", err
	}
	return service.generateToken(UUidString)
}

func (service *UserService) GetAccessToken(refreshToken string) (string, error) {
	userId, err := service.repository.GetUUid(refreshToken)
	if err != nil {
		return "", err
	}
	return service.generateToken(userId)
}

func (service *UserService) generateToken(userId string) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, models.MyJwtClaims{UserId: userId, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().UTC().Add(time.Minute * 5).Unix()}}).SignedString([]byte(service.config.JWTString))
	if err != nil {
		return "", errors.New("internal error")
	}
	return token, nil
}

func (service *UserService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) { return []byte(service.config.JWTString), nil })
	if err != nil {
		return "", err
	}
	mapClaims, okClaims := token.Claims.(jwt.MapClaims)
	user_id, okMap := mapClaims["user_id"].(string)
	if okClaims && okMap {
		return user_id, nil
	}
	return "", errors.New("user isn't verified by server")
}

func (service *UserService) hashPassword(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return errors.New("internal error")
	}
	user.Password = string(hashedPassword)
	return nil
}

func NewUserService(repository repository.Repository, config *config.Config, validator *validator.Validate) *UserService {
	return &UserService{
		repository: repository,
		config:     config,
		validator:  validator,
	}
}
