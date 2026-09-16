package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Expense struct {
	Amount      float64
	Category    string
	Description string
	Currency    string
}

func createExpense(amount float64, category, description, currency string) (Expense, error) {
	if amount <= 0 {
		return Expense{}, fmt.Errorf("сумма должна быть больше нуля")
	}

	if category == "" {
		return Expense{}, fmt.Errorf("категория не может быть пустой")
	}

	if description == "" {
		return Expense{}, fmt.Errorf("описание не может быть пустым")
	}

	if currency == "" {
		return Expense{}, fmt.Errorf("валюта не может быть пустой")
	}

	return Expense{
		Amount:      amount,
		Category:    category,
		Description: description,
		Currency:    currency,
	}, nil
}

func (e Expense) IsValid() bool {
	return e.Amount > 0 && e.Category != ""
}

func (e *Expense) UpdateDescription(description string) {
	e.Description = description
}

func (e Expense) Format() string {
	return fmt.Sprintf(
		"%.2f %s | %s | %s",
		e.Amount,
		e.Currency,
		e.Category,
		e.Description,
	)
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не установлен")
	}

	bot, err := tgbotapi.NewBotAPI(token)
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

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Привет! Я бот для учёта расходов.",
				)

				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}