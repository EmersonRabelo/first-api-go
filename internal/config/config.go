package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SettingProvider interface {
	GetEnvironment() string
	GetServer() Server
	GetDatabase() Database
	IsProd() bool
	IsTest() bool
	IsLocal() bool
}

type Database struct {
	Host    string
	Port    string
	User    string
	Pwd     string
	Name    string
	SSLMode string
}

type Server struct {
	Port string
}

type Setting struct {
	Environment string
	Server      Server
	Database    Database
}

var AppSetting SettingProvider

func GetSetting() SettingProvider {
	AppSetting = loadSetting()

	return AppSetting
}

func loadSetting() *Setting {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Arquivo %s não encontrado, tentando .env padrão", envFile)
		_ = godotenv.Load()
		// se não achar nada, usa do sistema mesmo
	}
	return &Setting{
		Environment: env,
		Server: Server{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
		},
		Database: Database{
			Host:    getEnvOrDefault("DB_HOST", "localhost"),
			Port:    getEnvOrDefault("DB_PORT", "5432"),
			User:    getEnvOrDefault("DB_USER", "postgres"),
			Pwd:     getEnvOrDefault("DB_PASSWORD", "1234"),
			Name:    getEnvOrDefault("DB_NAME", "go-api"),
			SSLMode: getEnvOrDefault("DB_SSL_MODE", "disable"),
		},
	}
}

func getEnvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func (s *Setting) GetEnvironment() string { return s.Environment }
func (s *Setting) GetServer() Server      { return s.Server }
func (s *Setting) GetDatabase() Database  { return s.Database }
func (s *Setting) IsProd() bool           { return s.Environment == "production" }
func (s *Setting) IsTest() bool           { return s.Environment == "test" }
func (s *Setting) IsLocal() bool          { return s.Environment == "local" || s.Environment == "development" }
