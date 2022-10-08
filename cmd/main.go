package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/logger"
	"github.com/Kin-dza-dzaa/userApi/internal/validation"
	"github.com/Kin-dza-dzaa/userApi/pkg/handlers"
	repository "github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
)

func main() {
	writer, err := logger.GetWriter()
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			log.Panic(err)
		}
	}()
	logger := logger.GetLogger(writer)
	config, err := config.ReadConfig(logger)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	pool, err := pgxpool.Connect(context.TODO(), config.DbUrl)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	validator, err := validation.InitValidators()
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	repository := repository.NewRepository(pool)
	service := service.NewService(repository, config, validator)
	handlers := handlers.NewHandlers(service, *config)
	srv := &http.Server{
		Addr: config.Adress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler: cors.Default().Handler(handlers.Router),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Panic().Msg(err.Error())
		}
	}()
	//  := make(chan bool)
	// <-
}
