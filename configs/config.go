package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string `mapstructure:"APP_NAME"`
	Port string `mapstructure:"APP_PORT"`
	JwtSecret string `mapstructure:"JWT_SECRET"`
	DBHost string `mapstructure:"DB_HOST"`
	GinMode string `mapstructure:"GIN_MODE"`
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPass string `mapstructure:"REDIS_PASS"`
	SMTPHost string `mapstructure:"SMTP_HOST"`
	SMTPPort string `mapstructure:"SMTP_PORT"`
	SMTPUser string `mapstructure:"SMTP_USER"`
	SMTPPass string `mapstructure:"SMTP_PASS"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
        log.Printf("Error reading config file: %s\n", err)
		return nil, err
    }

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

