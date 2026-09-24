package config

import (
	"log"
	"os"
	"strconv"
	"strings"
    "github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DatabaseURL   string
	AllowedUserIDs  map[int64]bool
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

	allowed := make(map[int64]bool)
	raw := os.Getenv("ALLOWED_USER_IDS")
	if raw != "" {
		for _, idStr := range strings.Split(raw, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				log.Printf("Некорректный ID в ALLOWED_USER_IDS: %s", idStr)
				continue
			}
			allowed[id] = true
		}
	}

	return Config{
		TelegramToken:  token,
		DatabaseURL:    dbURL,
		AllowedUserIDs: allowed,
	}
}