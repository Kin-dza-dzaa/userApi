package apierror

import (
	"errors"
	"net/http"

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
			switch err := err.(type) {

			case *ErrorStruct:
				w.WriteHeader(err.Code)
				w.Write(err.Marshal())

			default:

				if errors.Is(err, puddle.ErrNotAvailable) {
					apierror.logger.Err(err).Msg(err.Error())
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write(NewErrorStruct("too many requests", "error", http.StatusTooManyRequests).Marshal())
					return
				}

				if errors.Is(err, puddle.ErrClosedPool) {
					apierror.logger.Err(err).Msg(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(NewErrorStruct("internal server error", "error", http.StatusInternalServerError).Marshal())
					return
				}
				
				apierror.logger.Err(err).Msg(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(NewErrorStruct("unexpected error", "error", http.StatusInternalServerError).Marshal())
			}

		}
	})
}

func NewApiError(logger *zerolog.Logger) *ApiError {
	return &ApiError{
		logger: logger,
	}
}
