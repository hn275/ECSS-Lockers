package internal

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	Env struct {
		HostEmail    string
		HostPassword string
		MailServer   string
		MailPort     int
		Domain       string
		CipherKey    []byte
		DatabaseURL  string
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to load .env")
	}

	Env.HostEmail = EnvVar("GMAIL_USER")
	Env.HostPassword = EnvVar("GMAIL_PASSWORD")
	Env.MailServer = "smtp.gmail.com"
	Env.MailPort = 587
	Env.Domain = EnvVar("ORIGIN")
	Env.CipherKey = []byte(EnvVar("CIPHER_KEY"))
	if len(Env.CipherKey) != 32 {
		logger.Warn("invalid value set for env $CIPHER_KEY, expected length of 32 bytes, got %d byte(s)",
			len(Env.CipherKey))
	}

	Env.DatabaseURL = fmt.Sprintf("%s?authToken=%s",
		EnvVar("TURSO_DATABASE_URL"),
		EnvVar("TURSO_AUTH_TOKEN"))
}

func EnvMust(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logger.Fatal("env variable not set: $%s", key)
	}
	return value
}

func EnvVar(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logger.Warn("env variable not set: $%s", key)
	}
	return value
}
