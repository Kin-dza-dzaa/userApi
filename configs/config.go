package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct{
	DbUrl     string		`mapstructure:"DB_URL"`
	Adress    string		`mapstructure:"ADRESS"`
	JWTString string		`mapstructure:"JWT_SECURE_STRING"`
	AllowedOrigins string	`mapstructure:"ALLOWED_ORIGINS"`
}

func ReadConfig(logger *zerolog.Logger) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs")
	config := new(Config)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(config); err != nil {
		logger.Panic().Msg(err.Error())
	}	
	return config, nil
}