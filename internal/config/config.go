package config

import (
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	MYSQLUserDSN               string = "MYSQL_USER_DSN"
	MYSQLTimeoutMs             string = "MYSQL_TIMEOUT_MS"
	MYSQLMaxOpenConn           string = "MYSQL_MAX_OPEN_CONN"
	MYSQLMaxConnLifeTimeSecond string = "MYSQL_MAX_CONN_LIFETIME_SECOND"
	MYSQLMaxIdleConn           string = "MYSQL_MAX_IDLE_CONN"

	APITimeout   string = "API_TIMEOUT"
	DefaultLimit string = "DEFAULT_LIMIT"
)

type Config struct {
	MYSQLUserDSN               string `validate:"required"`
	MYSQLTimeoutMs             int    `validate:"required"`
	MYSQLMaxOpenConn           int    `validate:"required"`
	MYSQLMaxConnLifeTimeSecond int    `validate:"required"`
	MYSQLMaxIdleConn           int    `validate:"required"`

	APITimeout   int `validate:"required"`
	DefaultLimit int `validate:"required"`
}

func New(validate *validator.Validate) Config {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("failed to load .env files, loading host env instead")
	}

	cfg := Config{
		MYSQLUserDSN:               os.Getenv(MYSQLUserDSN),
		MYSQLTimeoutMs:             getEnvInt(MYSQLTimeoutMs, os.Getenv(MYSQLTimeoutMs)),
		MYSQLMaxOpenConn:           getEnvInt(MYSQLMaxOpenConn, os.Getenv(MYSQLMaxOpenConn)),
		MYSQLMaxConnLifeTimeSecond: getEnvInt(MYSQLMaxConnLifeTimeSecond, os.Getenv(MYSQLMaxConnLifeTimeSecond)),
		MYSQLMaxIdleConn:           getEnvInt(MYSQLMaxIdleConn, os.Getenv(MYSQLMaxIdleConn)),

		APITimeout:   getEnvInt(APITimeout, os.Getenv(APITimeout)),
		DefaultLimit: getEnvInt(DefaultLimit, os.Getenv(DefaultLimit)),
	}

	if err := validate.Struct(cfg); err != nil {
		log.Panic().Err(err).Msg("failed to validate config")
	}

	return cfg
}

// convert env to int
func getEnvInt(env string, value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.
			Err(err).
			Stack().
			Str("env", env).
			Str("value", value).
			Msg("failed to convert string to int")
	}
	return i
}
