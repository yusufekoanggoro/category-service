package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigProvider interface {
	GetHTTPPort() string
	GetHTTPHost() string

	GetBookGRPCHost() string
	GetBookGRPCPort() string

	GetDBHost() string
	GetDBPort() string
	GetDBUser() string
	GetDBPassword() string
	GetDBName() string
	GetSSLMode() string
}

type EnvConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string

	HTTPHost string
	HTTPPort string

	BookGRPCHost string
	BookGRPCPort string
}

func (e *EnvConfig) GetHTTPHost() string { return e.HTTPHost }
func (e *EnvConfig) GetHTTPPort() string { return e.HTTPPort }

func (e *EnvConfig) GetBookGRPCHost() string { return e.BookGRPCHost }
func (e *EnvConfig) GetBookGRPCPort() string { return e.BookGRPCPort }

func (e *EnvConfig) GetDBHost() string     { return e.DBHost }
func (e *EnvConfig) GetDBPort() string     { return e.DBPort }
func (e *EnvConfig) GetDBUser() string     { return e.DBUser }
func (e *EnvConfig) GetDBPassword() string { return e.DBPassword }
func (e *EnvConfig) GetDBName() string     { return e.DBName }
func (e *EnvConfig) GetSSLMode() string    { return e.SSLMode }

func LoadConfig() ConfigProvider {
	err := godotenv.Load()
	if err != nil {
		log.Println("Gagal membaca file .env, menggunakan environment variables yang tersedia")
	}

	return &EnvConfig{
		HTTPHost: os.Getenv("HTTP_HOST"),
		HTTPPort: os.Getenv("HTTP_PORT"),

		BookGRPCHost: os.Getenv("BOOK_GRPC_HOST"),
		BookGRPCPort: os.Getenv("BOOK_GRPC_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		SSLMode:    os.Getenv("DB_SSLMODE"),
	}
}
