package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Handlers struct {
	Router  *mux.Router
	Cors    *cors.Cors
	Service service.Service
}

func (handlers *Handlers) SignUpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "wrong input"})
			return
		}
		if err := handlers.Service.SignUpUser(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "wrong input"})
			return
		}
		jwt, err := handlers.Service.SignInUser(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Refresh-token",
			Value:    user.RefreshToken,
			Expires:  time.Now().UTC().AddDate(0, 6, 0),
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "token": jwt})
	})
}

func (handlers *Handlers) VerifyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var user models.User
		vars := mux.Vars(r)
		user.VerificationCode = vars["code"]
		jwt, err := handlers.Service.VerifyUser(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Refresh-token",
			Value:    user.RefreshToken,
			Expires:  time.Now().UTC().AddDate(0, 6, 0),
			HttpOnly: true,
			Secure:   false, // set true on realese
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "token": jwt})
	})
}

func (handlers *Handlers) GetTokenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		cookie, err := r.Cookie("Refresh-token")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "cooike isn't present"})
			return
		}
		jwt, err := handlers.Service.GetAccessToken(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "token": jwt})
	})
}

func NewHandlers(service service.Service, config config.Config) *Handlers {
	handlers := new(Handlers)
	handlers.Service = service
	handlers.Router = mux.NewRouter().Host(config.Adress).Subrouter()
	handlers.Cors = cors.New(cors.Options{
		AllowedOrigins: strings.Split(config.AllowedOrigins, ","),
		AllowedHeaders: []string{"User-Agent", "Content-type"},
		MaxAge:         5,
		AllowedMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	})

	user := handlers.Router.Handle("/user", handlers.SignUpHandler()).Methods("POST").Subrouter()
	user.Handle("/auth", handlers.SignInHandler()).Methods("POST")
	user.Handle("/token", handlers.GetTokenHandler()).Methods("GET")
	user.Handle("/verify/{code:[.{16}]}", handlers.VerifyHandler()).Methods("POST")
	return handlers
}
