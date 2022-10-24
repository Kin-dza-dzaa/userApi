package apierror

import (
	"errors"
	"net/http"

	"github.com/Kin-dza-dzaa/userApi/internal/dto"
	repository "github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/jackc/puddle"
	"github.com/rs/zerolog"
)

type UserHandler func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	logger *zerolog.Logger
}

func (apierror *ApiError) ErrorMiddleWare(next UserHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			apierror.logger.Err(err).Msg(err.Error())
			switch {
			case errors.Is(err, puddle.ErrClosedPool):
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(NewErrorStruct("internal server error", "error").Marshal())

			case errors.Is(err, puddle.ErrNotAvailable):
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write(NewErrorStruct("too many requests", "error").Marshal())

			case errors.Is(err, dto.ErrInvalidCredentials):
				w.WriteHeader(http.StatusBadRequest)
				w.Write(NewErrorStruct(dto.ErrInvalidCredentials.Error(), "error").Marshal())

			case errors.Is(err, service.ErrWrongPassowrd):
				w.WriteHeader(http.StatusBadRequest)
				w.Write(NewErrorStruct(service.ErrWrongPassowrd.Error(), "error").Marshal())

			case errors.Is(err, repository.ErrUserDoesntExists):
				w.WriteHeader(http.StatusBadRequest)
				w.Write(NewErrorStruct(repository.ErrUserDoesntExists.Error(), "error").Marshal())

			case errors.Is(err, repository.ErrWrongEmail):
				w.WriteHeader(http.StatusBadRequest)
				w.Write(NewErrorStruct(repository.ErrWrongEmail.Error(), "error").Marshal())

			case errors.Is(err, repository.ErrWrongVerificationCode):
				w.WriteHeader(http.StatusBadRequest)
				w.Write(NewErrorStruct(repository.ErrWrongVerificationCode.Error(), "error").Marshal())
				
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(NewErrorStruct("unexpected error", "error").Marshal())
			}
		}
	})
}

func NewApiError(logger *zerolog.Logger) *ApiError {
	return &ApiError{
		logger: logger,
	}
}
