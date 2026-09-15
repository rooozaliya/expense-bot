package main

import (
	"errors"
	"fmt"
)

type Expense struct {
	Amount float64
	Category string 
	Description string
}

func createExpense(amount float64, category string, description string) (Expense, error) {
	if(amount<= 0) {
		return Expense{}, errors.New("сумма должна быть больше нуля")
	}
	return Expense{
		Amount: amount,
		Category: category,
		Description: description,
	}, nil
}


func main() {
	expense, err := createExpense(-450, "food", "dinner")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}


	fmt.Println("Сумма:", expense.Amount)
	fmt.Println("Категория:", expense.Category)
	fmt.Println("Описание:", expense.Description)
}

