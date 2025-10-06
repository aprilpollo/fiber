package core

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App     *App
		JWT     *JWT
		Database *Database
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
		URI_Script      string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime int
	}
)

func GetConfig() (*Config, error) {
	if os.Getenv("APP_MODE") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}
	app := &App{
		AppName:                  os.Getenv("APP_NAME"),
		AppVersion:               os.Getenv("APP_VERSION"),
		ApiPort:                  os.Getenv("API_PORT"),
		ShutdownTimeout:          getEnvAsUint("API_SHUTDOWN_TIMEOUT_SECONDS", 30),
		AllowedCredentialOrigins: os.Getenv("ALLOWED_CREDENTIAL_ORIGINS"),
		LogLevel:                 os.Getenv("LOG_LEVEL"),
		Development:              os.Getenv("APP_MODE") == "development",
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
		URI_Script:      os.Getenv("POSTGRE_URI_SCRIPT"),
		MaxIdleConns:    getEnvAsInt("POSTGRE_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvAsInt("POSTGRE_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime: getEnvAsInt("POSTGRE_CONN_MAX_LIFETIME", 0),
	}

	return &Config{
		App:     app,
		JWT:     jwt,
		Database: database,
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
