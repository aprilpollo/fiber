package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App      *App
		JWT      *JWT
		Database *Database
		Redis    *Redis
		Oauth    *Oauth
		Minios3  *Minios3
		EmailAPI *EmailAPI
		Webhook  *Webhook
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
	}

	App struct {
		AppName                  string
		AppVersion               string
		ApiPort                  string
		ShutdownTimeout          uint
		AllowedCredentialOrigins string
		LogLevel                 string
		Development              bool
		TimeZone                 string
	}

	JWT struct {
		SecretKey          string
		JwtExpireDaysCount int
		Issuer             string
		Subject            string
		SigningMethod      string
	}

	Database struct {
		URI             string
		URI_MIGRATION   string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime int
	}

	Oauth struct {
		GoogleProvider GoogleProvider
	}

	GoogleProvider struct {
		ClientID string
	}

	Minios3 struct {
		Endpoint       string
		EndpointPublic string
		AccessKey      string
		SecretKey      string
		Bucket         string
		ExpireDays     int
		UseSSL         bool
	}

	EmailAPI struct {
		URL string
	}

	Webhook struct {
		URL    string
		Expire int //minutes
	}
)

func GetConfig() (*Config, error) {
	_ = godotenv.Load()
	app := &App{
		AppName:                  os.Getenv("APP_NAME"),
		AppVersion:               os.Getenv("APP_VERSION"),
		ApiPort:                  os.Getenv("API_PORT"),
		ShutdownTimeout:          getEnvAsUint("API_SHUTDOWN_TIMEOUT_SECONDS", 30),
		AllowedCredentialOrigins: os.Getenv("ALLOWED_CREDENTIAL_ORIGINS"),
		TimeZone:                 os.Getenv("TIME_ZONE"),
	}

	jwt := &JWT{
		SecretKey:          os.Getenv("JWT_SECRET_KEY"),
		JwtExpireDaysCount: getEnvAsInt("JWT_EXPIRE_DAYS_COUNT", 7),
		Issuer:             os.Getenv("JWT_ISSUER"),
		Subject:            os.Getenv("JWT_SUBJECT"),
		SigningMethod:      os.Getenv("JWT_SIGNING_METHOD"),
	}

	database := &Database{
		URI:             os.Getenv("POSTGRE_URI"),
		URI_MIGRATION:   os.Getenv("POSTGRE_URI_MIGRATION"),
		MaxIdleConns:    getEnvAsInt("POSTGRE_MAX_IDLE_CONNS", 2),
		MaxOpenConns:    getEnvAsInt("POSTGRE_MAX_OPEN_CONNS", 10),
		ConnMaxLifetime: getEnvAsInt("POSTGRE_CONN_MAX_LIFETIME", 300),
	}

	redis := &Redis{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       getEnvAsInt("REDIS_DB", 0),
	}

	oauth := &Oauth{
		GoogleProvider: GoogleProvider{
			ClientID: os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
		},
	}

	minios3 := &Minios3{
		Endpoint:       os.Getenv("MINIO_ENDPOINT"),
		EndpointPublic: os.Getenv("MINIO_ENDPOINT_PUBLIC"),
		AccessKey:      os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey:      os.Getenv("MINIO_SECRET_KEY"),
		Bucket:         os.Getenv("MINIO_BUCKET"),
		ExpireDays:     getEnvAsInt("MINIO_EXPIRE_DAY", 7),
		UseSSL:         os.Getenv("MINIO_USE_SSL") == "true",
	}

	emailAPI := &EmailAPI{
		URL: os.Getenv("EMAIL_API_URL"),
	}

	webhook := &Webhook{
		URL:    os.Getenv("WEBHOOK_URL"),
		Expire: getEnvAsInt("WEBHOOK_EXPIRE", 10),
	}

	return &Config{
		App:      app,
		JWT:      jwt,
		Database: database,
		Redis:    redis,
		Oauth:    oauth,
		Minios3:  minios3,
		EmailAPI: emailAPI,
		Webhook:  webhook,
	}, nil
}

func getEnvAsUint(key string, defaultValue uint) uint {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint(parsed)
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
