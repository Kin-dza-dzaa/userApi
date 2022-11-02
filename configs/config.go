package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	DbUrl            string `mapstructure:"DB_URL"`
	Adress           string `mapstructure:"ADRESS"`
	JWTString        string `mapstructure:"JWT_SECURE_STRING"`
	AllowedOrigins   string `mapstructure:"ALLOWED_ORIGINS"`
	SmtpUserName     string `mapstructure:"SMTP_USER_NAME"`
	SmtpPassword     string `mapstructure:"SMTP_PASSWORD"`
	TemplateLocation string `mapstructure:"TEMPLATE_LOCATION"`
	Secure           bool   `mapstructure:"SECURE"`
	LocalDbUrlTest   string `mapstructure:"LOCAL_DB_URL_TEST"`
	AllowCredentials bool   `mapstructure:"ALLOW_CREDENTIALS"`
	SpaUrl           string `mapstructure:"SPA_URL"`
}

func ReadConfig(logger *zerolog.Logger) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./../../configs") // for test
	viper.AddConfigPath("./../configs")    // for test
	config := new(Config)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		logger.Panic().Msg(err.Error())
	}
	return config, nil
}
