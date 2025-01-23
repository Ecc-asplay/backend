package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment         string        `mapstructure:"ENVIRONMENT"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	FrontAddress        []string      `mapstructure:"FRONT_ADDRESS"`
	AllowHeaders        []string      `mapstructure:"ALLOW_HEADERS"`
	RedisAddress        string        `mapstructure:"REDIS_ADDRESS"`
	RedisPassword       string        `mapstructure:"REDIS_PASSWORD"`
	MigrationURL        string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	SmtpHost            string        `mapstructure:"SMTP_HOST"`
	SmtpPort            int           `mapstructure:"SMTP_PORT"`
	SmtpUser            string        `mapstructure:"SMTP_USER"`
	SmtpPassword        string        `mapstructure:"SMTP_PASSWORD"`
	SmtpFromAddress     string        `mapstructure:"SMTP_FROM_ADDRESS"`
	SmtpFromName        string        `mapstructure:"SMTP_FROM_NAME"`
	AwsAccessKey        string        `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey        string        `mapstructure:"AWS_SECRET_KEY"`
	AwsRegion           string        `mapstructure:"AWS_REGION"`
	BucketName          string        `mapstructure:"BUCKET_NAME"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
