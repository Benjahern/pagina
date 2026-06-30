package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	envPort        = "PORT"
	envGinMode     = "GIN_MODE"
	envCORSOrigins = "CORS_ALLOWED_ORIGINS"
	envSeed        = "SEED"
	envLogLevel    = "LOG_LEVEL"

	envDBHost     = "DB_HOST"
	envDBPort     = "DB_PORT"
	envDBUser     = "DB_USER"
	envDBPassword = "DB_PASSWORD"
	envDBName     = "DB_NAME"

	envJWTSecret        = "JWT_SECRET"
	envJWTExpirationHrs = "JWT_EXPIRATION_HOURS"
	envBcryptCost       = "BCRYPT_COST"

	minJWTSecretBytes = 32
	minBcryptCost     = 4
	maxBcryptCost     = 31
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port        string
	GinMode     string
	CORSAllowed []string
	Seed        bool
	LogLevel    string
}

type DatabaseConfig struct {
	Host, User, Password, Name string
	Port                       int
}

type JWTConfig struct {
	Secret          []byte
	ExpirationHours int
	BcryptCost      int
}

// Load reads environment variables (optionally from .env), validates required
// values, and returns a typed Config. In production, .env is absent and env vars
// come from the orchestrator — that case is handled silently by godotenv.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("config: loading .env: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv(envPort, "8080"),
			GinMode:     getEnv(envGinMode, "debug"),
			CORSAllowed: getEnvList(envCORSOrigins, ",", []string{"http://localhost:3000"}),
			LogLevel:    getEnv(envLogLevel, "info"),
		},
	}

	seed, err := getEnvBool(envSeed, false)
	if err != nil {
		return nil, err
	}
	cfg.Server.Seed = seed

	dbHost, err := getEnvRequired(envDBHost)
	if err != nil {
		return nil, err
	}
	cfg.Database.Host = dbHost

	dbUser, err := getEnvRequired(envDBUser)
	if err != nil {
		return nil, err
	}
	cfg.Database.User = dbUser

	dbPassword, err := getEnvRequired(envDBPassword)
	if err != nil {
		return nil, err
	}
	cfg.Database.Password = dbPassword

	dbName, err := getEnvRequired(envDBName)
	if err != nil {
		return nil, err
	}
	cfg.Database.Name = dbName

	dbPort, err := getEnvIntRequired(envDBPort)
	if err != nil {
		return nil, err
	}
	cfg.Database.Port = dbPort

	jwtSecret, err := getEnvRequired(envJWTSecret)
	if err != nil {
		return nil, err
	}
	if len(jwtSecret) < minJWTSecretBytes {
		return nil, fmt.Errorf("config: %s must be at least %d bytes (got %d)", envJWTSecret, minJWTSecretBytes, len(jwtSecret))
	}
	cfg.JWT.Secret = []byte(jwtSecret)

	jwtExp, err := getEnvInt(envJWTExpirationHrs, 24)
	if err != nil {
		return nil, err
	}
	cfg.JWT.ExpirationHours = jwtExp

	bcryptCost, err := getEnvInt(envBcryptCost, 12)
	if err != nil {
		return nil, err
	}
	if bcryptCost < minBcryptCost || bcryptCost > maxBcryptCost {
		return nil, fmt.Errorf("config: %s must be between %d and %d (got %d)", envBcryptCost, minBcryptCost, maxBcryptCost, bcryptCost)
	}
	cfg.JWT.BcryptCost = bcryptCost

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return defaultValue
}

func getEnvRequired(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return "", fmt.Errorf("config: required env var %s is not set", key)
	}
	return v, nil
}

func getEnvInt(key string, defaultValue int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return defaultValue, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer (got %q)", key, v)
	}
	return n, nil
}

func getEnvIntRequired(key string) (int, error) {
	v, err := getEnvRequired(key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer (got %q)", key, v)
	}
	return n, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return defaultValue, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("config: %s must be a boolean (got %q)", key, v)
	}
	return b, nil
}

func getEnvList(key, sep string, defaultValue []string) []string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return defaultValue
	}
	parts := strings.Split(v, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return defaultValue
	}
	return out
}
