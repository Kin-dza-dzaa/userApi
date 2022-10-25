package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/apierror"
	"github.com/Kin-dza-dzaa/userApi/pkg/handlers"
	repository "github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(zerolog.SyncWriter(os.Stdout)).With().Timestamp().Caller().Logger()
	config, err := config.ReadConfig(&logger)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	pool, err := pgxpool.Connect(context.TODO(), config.DbUrl)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	ApiError := apierror.NewApiError(&logger)
	repository := repository.NewRepository(pool)
	service := service.NewService(repository, config)
	MyHandlers := handlers.NewHandlers(service, config, ApiError)
	srv := &http.Server{
		Addr:         config.Adress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      MyHandlers.Cors.Handler(MyHandlers.Router),
	}
	go func() {
		logger.Info().Msg(fmt.Sprintf("Staring server userapi at %v", config.Adress))
		if err := srv.ListenAndServe(); err != nil {
			logger.Panic().Msg(err.Error())
		}
	}()
	<-handlers.StopHTTPServerChan
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}
