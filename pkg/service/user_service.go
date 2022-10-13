package service

import (
	"bytes"
	"errors"
	"html/template"
	"time"

	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	repository "github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/dchest/uniuri"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type UserService struct {
	repository repository.Repository
	config     *config.Config
	validator  *validator.Validate
	logger 	   *zerolog.Logger
}

func (service *UserService) SignUpUser(user *models.User) error {
	if err := service.validator.Struct(user); err != nil {
		return errors.New("invalid credentials")
	}
	if err := service.hashPassword(user); err != nil {
		return err
	}
	user.UserId = uuid.New()
	user.RegistrationTime = time.Now().UTC()
	user.VerificationCode = uniuri.New()
	user.Verified = false
	
	exists, err := service.repository.IfUnverifiedUserExists(user)
	if err != nil {
		return err
	}

	if exists {
		if err := service.repository.UpdateCredentials(user); err != nil {
			return err
		}
	} else {
		if err := service.repository.AddUser(user); err != nil {
			return err
		}
	}

	if err := service.sendVerificationCode(user); err != nil {
		return err
	}
	return nil
}

func (service *UserService) SignInUser(user *models.User) error {
	if err := service.validator.Struct(user); err != nil {
		return errors.New("invalid credentials")
	}
	DbUser, err := service.repository.GetVerifiedUser(user)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(DbUser.Password), []byte(user.Password)); err != nil {
		return errors.New("invalid password")
	}
	if time.Now().After(DbUser.ExpirationTime) {
		user.RefreshToken = uniuri.NewLen(512)
		user.ExpirationTime = time.Now().UTC().AddDate(0, 6, 0)
		if err := service.repository.UpdateRefreshToken(user); err != nil {
			return err
		}
	} else {
		user.RefreshToken = DbUser.RefreshToken
		user.ExpirationTime = DbUser.ExpirationTime
	}
	return service.generateToken(user)
}

func (service *UserService) VerifyUser(user *models.User) error {
	user.RefreshToken = uniuri.NewLen(512)
	user.ExpirationTime = time.Now().UTC().AddDate(0, 6, 0)
	if err := service.repository.VerifyUser(user); err != nil {
		return err
	}
	err := service.repository.GetUUid(user)
	if err != nil {
		return err
	}
	return service.generateToken(user)
}

func (service *UserService) GetAccessToken(user *models.User) error {
	err := service.repository.GetUUid(user)
	if err != nil {
		return err
	}
	return service.generateToken(user)
}

func (service *UserService) generateToken(user *models.User) error {
	user.CsrfToken = uniuri.NewLen(32)
	jwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, models.MyJwtClaims{UserId: user.UserId.String(), XCSRFToken: user.CsrfToken, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().UTC().Add(time.Minute * 5).Unix()}}).SignedString([]byte(service.config.JWTString))
	if err != nil {
		service.logger.Error().Msg(err.Error())
		return errors.New("internal error")
	}
	user.Jwt = jwt
	return nil
}

func (service *UserService) hashPassword(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		service.logger.Error().Msg(err.Error())
		return errors.New("internal error")
	}
	user.Password = string(hashedPassword)
	return nil
}

func (service *UserService) sendVerificationCode(user *models.User) error {
	t, err := template.ParseFiles("./internal/templates/response_template.html")
	if err != nil {
		service.logger.Error().Msg(err.Error())
		return errors.New("internal error")
	}
	body := new(bytes.Buffer)
	err = t.Execute(body, map[string]string{"VerificationCode": user.VerificationCode, "UserName": user.UserName})
	if err != nil {
		service.logger.Error().Msg(err.Error())
		return errors.New("internal error")
	}
	d := gomail.NewDialer("smtp.gmail.com", 465, service.config.SmtpUserName, service.config.SmtpPassword)
	msg := gomail.NewMessage()
	msg.SetHeader("From", service.config.SmtpUserName)
    msg.SetHeader("To", user.Email)
    msg.SetHeader("Subject", "Verification on WordDict")
    msg.SetBody("text/html", body.String())
	if err := d.DialAndSend(msg); err != nil {
		service.logger.Error().Msg(err.Error())
		return errors.New("can't send verification email :(, try 5 min later")
	}
	return nil
}

func NewUserService(repository repository.Repository, config *config.Config, validator *validator.Validate, logger *zerolog.Logger) *UserService {
	return &UserService{
		repository: repository,
		config:     config,
		validator:  validator,
		logger: logger,
	}
}
