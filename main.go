package main

import "fmt"

type Expense struct {
	Amount float64
	Category string 
	Description string
}

func createExpense(amount float64, category string, description string) *Expense {
	return &Expense{
		Amount: amount,
		Category: category,
		Description: description,
	}
}


func main() {
	expense:= createExpense(450, "food", "dinner")


	fmt.Println(expense)
	fmt.Println("Сумма:", expense.Amount)
	fmt.Println("Категория:", expense.Category)
	fmt.Println("Описание:", expense.Description)
}

