package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/apierror"
	"github.com/Kin-dza-dzaa/userApi/internal/dto"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var StopHTTPServerChan = make(chan bool)

var (
	ErrCookieNotPresent = errors.New("cookie not present")
)

type Handlers struct {
	Router  *mux.Router
	Cors    *cors.Cors
	Service service.Service
	ApiError *apierror.ApiError 
	Config *config.Config
}

func (handlers *Handlers) SignUpHandler() apierror.UserHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-type", "application/json")
		var userDto dto.UserSignUpDto
		if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
			return apierror.NewErrorStruct(`unmarshal failed, expected object{"email":"string", "password":"string", "user_name": "string"}`, "error", http.StatusBadRequest)
		}
		User, err := userDto.IntoUser()
		if err != nil {
			return err
		}
		if err := handlers.Service.SignUpUser(r.Context(), User); err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "email was sent", "code": 200})
		return nil
	}
}

func (handlers *Handlers) SignInHandler() apierror.UserHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-type", "application/json")
		var userDto dto.UserSignInDto
		if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
			return apierror.NewErrorStruct(`unmarshal failed, expected object{"email":"string", "password":"string"}`, "error", http.StatusBadRequest)
		}
		User, err := userDto.IntoUser()
		if err != nil {
			return err
		}
		if err := handlers.Service.SignInUser(r.Context(), User); err != nil {
			return err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Refresh-token",
			Value:    User.RefreshToken,
			Expires:  User.ExpirationTime,
			HttpOnly: true,
			Secure:   handlers.Config.Secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "Access-token",
			Value:    User.Jwt,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   handlers.Config.Secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.Header().Set("X-CSRF-Token", User.CsrfToken)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "code": 200})
		return nil
	}
}

func (handlers *Handlers) VerifyHandler() apierror.UserHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-type", "application/json")
		var user models.User
		vars := mux.Vars(r)
		user.VerificationCode = vars["code"]
		if err := handlers.Service.VerifyUser(r.Context(), &user); err != nil {
			return err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Refresh-token",
			Value:    user.RefreshToken,
			Expires:  user.ExpirationTime,
			HttpOnly: true,
			Secure:   handlers.Config.Secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "Access-token",
			Value:    user.Jwt,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   handlers.Config.Secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.Header().Set("X-CSRF-Token", user.CsrfToken)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "code": 200})
		return nil
	}
}

func (handlers *Handlers) GetTokenHandler() apierror.UserHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-type", "application/json")
		cookie, err := r.Cookie("Refresh-token")
		if err != nil {
			return apierror.NewErrorStruct(ErrCookieNotPresent.Error(), "error", http.StatusBadRequest)
		}
		user := new(models.User)
		user.RefreshToken = cookie.Value
		if err := handlers.Service.GetAccessToken(r.Context(), user); err != nil {
			return err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Access-token",
			Value:    user.Jwt,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   handlers.Config.Secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.Header().Set("X-CSRF-Token", user.CsrfToken)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "code": 200})
		return nil
	}
}

func (handlers *Handlers) LogOutHandler() apierror.UserHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:     "Refresh-token",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   handlers.Config.Secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "Access-token",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   handlers.Config.Secure,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "code": 200})
		return nil
	}
}

func NewHandlers(service service.Service, config *config.Config, ApiError *apierror.ApiError) *Handlers {
	handlers := new(Handlers)
	handlers.ApiError = ApiError
	handlers.Service = service
	handlers.Config = config
	handlers.Router = mux.NewRouter()
	handlers.Cors = cors.New(cors.Options{
		AllowedOrigins: strings.Split(config.AllowedOrigins, ","),
		AllowedHeaders: []string{"User-Agent", "Content-type"},
		ExposedHeaders: []string{"X-Csrf-Token"},
		AllowCredentials: config.AllowCredentials,
		MaxAge:         5,
		AllowedMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	})
	handlers.Router.Handle("/user", handlers.ApiError.ErrorMiddleWare(handlers.SignUpHandler())).Methods("POST").Schemes("http")
	user := handlers.Router.PathPrefix("/user").Subrouter()
	user.Handle("/auth", handlers.ApiError.ErrorMiddleWare(handlers.SignInHandler())).Methods("POST").Schemes("http")
	user.Handle("/token", handlers.ApiError.ErrorMiddleWare(handlers.GetTokenHandler())).Methods("GET").Schemes("http")
	user.Handle("/verify/{code:.{16}}", handlers.ApiError.ErrorMiddleWare(handlers.VerifyHandler())).Methods("POST").Schemes("http")
	user.Handle("/logout", handlers.ApiError.ErrorMiddleWare(handlers.LogOutHandler())).Methods("GET").Schemes("http")
	return handlers
}
