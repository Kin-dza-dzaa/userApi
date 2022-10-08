package logger

import (
	"os"
	"github.com/rs/zerolog"
)

func GetWriter() (*os.File, error) {
	return os.OpenFile("./internal/logger/logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
}

func GetLogger(writer *os.File) *zerolog.Logger{
	logger := zerolog.New(zerolog.SyncWriter(writer)).With().Timestamp().Caller().Logger()
	return &logger
}
