package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/internal/validation"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

var StopHTTPServerChan = make(chan bool)


type Handlers struct {
	Router  *mux.Router
	Cors    *cors.Cors
	Service service.Service
	Logger *zerolog.Logger
}

func (handlers *Handlers) SignUpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(validation.ParseError(err))
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "wrong input"})
			return
		}
		if err := handlers.Service.SignUpUser(&user); err != nil {
			w.WriteHeader(validation.ParseError(err))
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "email was sent"})
	})
}

func (handlers *Handlers) SignInHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(validation.ParseError(err))
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "wrong input"})
			return
		}
		err := handlers.Service.SignInUser(&user)
		if err != nil {
			w.WriteHeader(validation.ParseError(err))
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Refresh-token",
			Value:    user.RefreshToken,
			Expires:  user.ExpirationTime,
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "Access-token",
			Value:    user.Jwt,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.Header().Set("X-CSRF-Token", user.CsrfToken)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok"})
	})
}

func (handlers *Handlers) VerifyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var user models.User
		vars := mux.Vars(r)
		user.VerificationCode = vars["code"]
		err := handlers.Service.VerifyUser(&user)
		if err != nil {
			w.WriteHeader(validation.ParseError(err))
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Refresh-token",
			Value:    user.RefreshToken,
			Expires:  user.ExpirationTime,
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "Access-token",
			Value:    user.Jwt,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.Header().Set("X-CSRF-Token", user.CsrfToken)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok"})
	})
}

func (handlers *Handlers) GetTokenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		cookie, err := r.Cookie("Refresh-token")
		if err != nil {
			w.WriteHeader(validation.ParseError(err))
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "cooike isn't present"})
			return
		}
		user := new(models.User)
		user.RefreshToken = cookie.Value
		err = handlers.Service.GetAccessToken(user)
		if err != nil {
			w.WriteHeader(validation.ParseError(err))
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Access-token",
			Value:    user.Jwt,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.Header().Set("X-CSRF-Token", user.CsrfToken)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok"})
	})
}

func (handlers *Handlers) LogOutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:     "Refresh-token",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "Access-token",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok"})
	})
}

func NewHandlers(service service.Service, config config.Config, Logger *zerolog.Logger) *Handlers {
	handlers := new(Handlers)
	handlers.Logger = Logger	
	handlers.Service = service
	handlers.Router = mux.NewRouter().Host(config.Adress).Subrouter()
	handlers.Cors = cors.New(cors.Options{
		AllowedOrigins: strings.Split(config.AllowedOrigins, ","),
		AllowedHeaders: []string{"User-Agent", "Content-type"},
		MaxAge:         5,
		AllowedMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	})

	handlers.Router.Handle("/user", handlers.SignUpHandler()).Methods("POST").Schemes("http")
	user := handlers.Router.PathPrefix("/user").Subrouter()
		user.Handle("/auth", handlers.SignInHandler()).Methods("POST").Schemes("http")
		user.Handle("/token", handlers.GetTokenHandler()).Methods("GET").Schemes("http")
		user.Handle("/verify/{code:.{16}}", handlers.VerifyHandler()).Methods("POST").Schemes("http")
		user.Handle("/logout", handlers.LogOutHandler()).Methods("GET").Schemes("http")
		
	return handlers
}
