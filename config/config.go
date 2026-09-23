package config

import (
	"log"
	"os"
    "github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DatabaseURL   string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, использую системные переменные")
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не установлен")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://expense_user:expense_pass@localhost:5432/expense_bot"
	}

	return Config{TelegramToken: token, DatabaseURL: dbURL}
}