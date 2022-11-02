package service

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"text/template"
	"time"

	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/apierror"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	repository "github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/dchest/uniuri"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

var (
	ErrWrongPassowrd = errors.New("wrong password")
)

type UserService struct {
	repository repository.Repository
	config     *config.Config
}

func (service *UserService) SignUpUser(ctx context.Context, user *models.User) error {
	if err := service.hashPassword(user); err != nil {
		return err
	}
	user.UserId = uuid.New()
	user.RegistrationTime = time.Now().UTC()
	user.VerificationCode = uniuri.New()
	user.Verified = false
	exists := new(bool)
	err := service.repository.IfUnverifiedUserExists(ctx, user, exists)
	if err != nil {
		return err
	}

	if *exists {
		if err := service.repository.UpdateCredentials(ctx, user); err != nil {
			return err
		}
	} else {
		if err := service.repository.AddUser(ctx, user); err != nil {
			return err
		}
	}

	if err := service.sendVerificationCode(user); err != nil {
		return err
	}

	return nil
}

func (service *UserService) SignInUser(ctx context.Context, user *models.User) error {
	DbUser, err := service.repository.GetVerifiedUser(ctx, user)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(DbUser.Password), []byte(user.Password)); err != nil {
		return apierror.NewErrorStruct(ErrWrongPassowrd.Error(), "error", http.StatusBadRequest)
	}
	if time.Now().After(DbUser.ExpirationTime) {
		user.UserId = DbUser.UserId
		user.RefreshToken = uniuri.NewLen(512)
		user.ExpirationTime = time.Now().UTC().AddDate(0, 6, 0)
		if err := service.repository.UpdateRefreshToken(ctx, user); err != nil {
			return err
		}
	} else {
		user.UserId = DbUser.UserId
		user.RefreshToken = DbUser.RefreshToken
		user.ExpirationTime = DbUser.ExpirationTime
	}
	return service.generateToken(user)
}

func (service *UserService) VerifyUser(ctx context.Context, user *models.User) error {
	user.RefreshToken = uniuri.NewLen(512)
	user.ExpirationTime = time.Now().UTC().AddDate(0, 6, 0)
	if err := service.repository.VerifyUser(ctx, user); err != nil {
		return err
	}
	err := service.repository.GetUUid(ctx, user)
	if err != nil {
		return err
	}
	return service.generateToken(user)
}

func (service *UserService) GetAccessToken(ctx context.Context, user *models.User) error {
	err := service.repository.GetUUid(ctx, user)
	if err != nil {
		return err
	}
	return service.generateToken(user)
}

func (service *UserService) generateToken(user *models.User) error {
	user.CsrfToken = uniuri.NewLen(32)
	jwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, models.MyJwtClaims{UserId: user.UserId.String(), XCSRFToken: user.CsrfToken, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().UTC().Add(time.Minute * 5).Unix()}}).SignedString([]byte(service.config.JWTString))
	if err != nil {
		return err
	}
	user.Jwt = jwt
	return nil
}

func (service *UserService) hashPassword(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (service *UserService) sendVerificationCode(user *models.User) error {
	t, err := template.ParseFiles(service.config.TemplateLocation)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = t.Execute(body, map[string]string{"VerificationCode": service.config.SpaUrl + "/verify/" + user.VerificationCode, "UserName": user.UserName})
	if err != nil {
		return err
	}
	d := gomail.NewDialer("smtp.gmail.com", 465, service.config.SmtpUserName, service.config.SmtpPassword)
	msg := gomail.NewMessage()
	msg.SetHeader("From", service.config.SmtpUserName)
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", "Verification on WordDict")
	msg.SetBody("text/html", body.String())
	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

func NewUserService(repository repository.Repository, config *config.Config) *UserService {
	return &UserService{
		repository: repository,
		config:     config,
	}
}
