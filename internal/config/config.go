package config

import (
	"os"
)

type (
	DBConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Dbname   string
	}
	App struct {
		Mode string
		Port string
	}
	Cache struct {
		TTL int
	}
	Config struct {
		DBConfig *DBConfig
		App      *App
		Cache    *Cache
		AdminMW  *AdminMW
	}
	AdminMW struct {
		Key string
	}
)

func New() *Config {
	return &Config{
		DBConfig: &DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Dbname:   os.Getenv("DB_NAME"),
		},
		Cache: &Cache{
			TTL: 3600,
		},
		AdminMW: &AdminMW{
			Key: os.Getenv("ADMIN_KEY"),
		},
		App: &App{
			Mode: os.Getenv("APP_MODE"),
			Port: os.Getenv("APP_PORT"),
		},
	}
}
