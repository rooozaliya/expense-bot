package main


import (
	"log"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"expense-bot/config"
	"expense-bot/handlers"
	"expense-bot/repository"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Подключение к БД установлено")

	repo := repository.NewExpenseRepository(db)
	h := handlers.New(repo)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Бот запущен: @%s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		h.HandleUpdate(bot, update)
	}

}




