package models

import (
	"fmt"
	"time"
)


type Expense struct {
	Amount      float64
	Category    string
	Description string
	Currency    string
	CreatedAt   time.Time
}

type UserState struct {
	Step        string  // "" | "waiting_amount" | "waiting_category" | "waiting_description"
	Amount      float64
	Category    string
}



func CreateExpense(amount float64, category, description, currency string) (Expense, error) {
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
