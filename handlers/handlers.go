package handlers

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"expense-bot/models"
	"expense-bot/repository"
)

type Handler struct {
	repo   *repository.ExpenseRepository
	states map[int64]*models.UserState
}

func New(repo *repository.ExpenseRepository) *Handler {
	return &Handler{
		repo:   repo,
		states: make(map[int64]*models.UserState),
	}
}


func (h *Handler) sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func (h *Handler) handleTextMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	state, exists := h.states[chatID]
	if !exists {
		h.sendMessage(bot, chatID, "Не понимаю. Используй /add чтобы добавить расход.")
		return
	}

	switch state.Step {
	case "waiting_amount":
		amount, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			h.sendMessage(bot, chatID, "Это не похоже на число. Введи сумму ещё раз, например: 350.50")
			return // остаёмся на том же шаге, ждём повторный ввод
		}
		if amount <= 0 {
			h.sendMessage(bot, chatID, "Сумма должна быть больше нуля. Попробуй ещё раз:")
			return
		}
		state.Amount = amount
		state.Step = "waiting_category"
		h.sendMessage(bot, chatID, "Принято. Теперь введи категорию (например: еда, транспорт):")

	case "waiting_category":
		category := strings.TrimSpace(text)
		if category == "" {
			h.sendMessage(bot, chatID, "Категория не может быть пустой. Введи ещё раз:")
			return
		}

		state.Category = category
		state.Step = "waiting_description"
		h.sendMessage(bot, chatID, "Отлично. Добавь описание (или отправь '-' чтобы пропустить):")

	case "waiting_description":
		description := strings.TrimSpace(text)
		if description == "-" {
			description = ""
		}

		expense, err := models.CreateExpense(state.Amount, state.Category, description, "RUB")
		if err != nil {
			h.sendMessage(bot, chatID, "Ошибка: "+err.Error())
			return
		}
		
		if err := h.repo.Add(chatID, expense); err != nil {
			log.Println("Ошибка сохранения расхода:", err)
			h.sendMessage(bot, chatID, "Не удалось сохранить расход, попробуй позже")
			return
		}
		
		h.sendMessage(bot, chatID, "Расход добавлен:\n"+expense.Format())
		delete(h.states, chatID)
	}
}


func (h *Handler) handleListCommand(bot *tgbotapi.BotAPI, chatID int64) {
	expenses, err := h.repo.GetAll(chatID)

	if err != nil {
		log.Println("Ошибка получения расходов:", err)
		h.sendMessage(bot, chatID, "Не удалось получить список расходов")
		return
	}

	if len(expenses) == 0 {
		h.sendMessage(bot, chatID, "У тебя пока нет расходов. Добавь через /add")
		return
	}

	var lines []string
	for _, expense := range expenses {
		lines = append(lines, expense.Format())
	}

	text := strings.Join(lines, "\n")
	h.sendMessage(bot, chatID, text)
}



func (h *Handler) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	chatID := update.Message.Chat.ID

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			h.sendMessage(bot, chatID, "Привет! Я бот для учёта расходов.\nКоманды: /add, /help")

		case "help":
			h.sendMessage(bot, chatID, "/add — добавить расход\n/help — эта справка")
		case "add":
			h.states[chatID] = &models.UserState{Step: "waiting_amount"}
			h.sendMessage(bot, chatID, "Введи сумму расхода:")
		case "list":
			h.handleListCommand(bot, chatID)
		}
		
		return
		
		}
	h.handleTextMessage(bot, chatID, update.Message.Text)
}