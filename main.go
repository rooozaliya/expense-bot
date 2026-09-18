package main


import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Expense struct {
	Amount      float64
	Category    string
	Description string
	Currency    string
	CreatedAt   time.Time
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
		CreatedAt:   time.Now(),
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
		"%.2f %s | %s | %s | %s",
		e.Amount,
		e.Currency,
		e.Category,
		e.Description,
		e.CreatedAt.Format("02.01.2006 15:04"),
	)
}

type UserState struct {
	Step        string  // "" | "waiting_amount" | "waiting_category" | "waiting_description"
	Amount      float64
	Category    string
}

// userStates хранит состояние каждого пользователя по его chat ID
var userStates = make(map[int64]*UserState)

var userExpenses = make(map[int64][]Expense)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Println("Файл .env не найден, использую системные переменные")
    }

    token := os.Getenv("TELEGRAM_BOT_TOKEN")
	log.Println(token)

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

		chatID := update.Message.Chat.ID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				sendMessage(bot, chatID, "Привет! Я бот для учёта расходов.\nКоманды: /add, /help")

			case "help":
				sendMessage(bot, chatID, "/add — добавить расход\n/help — эта справка")
			case "add":
				userStates[chatID] = &UserState{Step: "waiting_amount"}
				sendMessage(bot, chatID, "Введи сумму расхода:")
			case "list":
				handleListCommand(bot, chatID)
			}
			
			continue
			
		}
	handleTextMessage(bot, chatID, update.Message.Text)

	}

}


func addExpense(chatID int64, expense Expense) {
	userExpenses[chatID] = append(userExpenses[chatID], expense)
}

func getExpenses(chatID int64) []Expense {
	return userExpenses[chatID] // если пусто — вернётся nil
}

func handleListCommand(bot *tgbotapi.BotAPI, chatID int64) {
	expenses := getExpenses(chatID)

	if len(expenses) == 0 {
		sendMessage(bot, chatID, "У тебя пока нет расходов. Добавь через /add")
		return
	}

	var lines []string
	for _, expense := range expenses {
		lines = append(lines, expense.Format())
	}

	text := strings.Join(lines, "\n")
	sendMessage(bot, chatID, text)
}

//отправка сообщений
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}

}

func handleTextMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	state, exists := userStates[chatID]
	if !exists {
		sendMessage(bot, chatID, "Не понимаю. Используй /add чтобы добавить расход.")
		return
	}

	switch state.Step {
	case "waiting_amount":
		amount, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			sendMessage(bot, chatID, "Это не похоже на число. Введи сумму ещё раз, например: 350.50")
			return // остаёмся на том же шаге, ждём повторный ввод
		}
		if amount <= 0 {
			sendMessage(bot, chatID, "Сумма должна быть больше нуля. Попробуй ещё раз:")
			return
		}
		state.Amount = amount
		state.Step = "waiting_category"
		sendMessage(bot, chatID, "Принято. Теперь введи категорию (например: еда, транспорт):")

	case "waiting_category":
		category := strings.TrimSpace(text)
		if category == "" {
			sendMessage(bot, chatID, "Категория не может быть пустой. Введи ещё раз:")
			return
		}

		state.Category = category
		state.Step = "waiting_description"
		sendMessage(bot, chatID, "Отлично. Добавь описание (или отправь '-' чтобы пропустить):")

	case "waiting_description":
		description := strings.TrimSpace(text)
		if description == "-" {
			description = ""
		}

		expense, err := createExpense(state.Amount, state.Category, description, "RUB")
		if err != nil {
			sendMessage(bot, chatID, "Ошибка: "+err.Error())
			return
		}
		addExpense(chatID, expense) 

		sendMessage(bot, chatID, "Расход добавлен:\n"+expense.Format())

		delete(userStates, chatID) // диалог завершён, сбрасываем состояние
	
	}
}
