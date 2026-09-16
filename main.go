package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Expense struct {
	Amount float64
	Category string 
	Description string
}

func createExpense(amount float64, category string, description string) (Expense, error) {
	if(amount<= 0) {
		return Expense{}, fmt.Errorf("сумма должна быть больше нуля")
	}


	if category == "" {
		return Expense{}, fmt.Errorf("категория не может быть пустой")
	}

	if description == "" {
		return Expense{}, fmt.Errorf("описание не может быть пустой")
	}

	return Expense{
		Amount: amount,
		Category: category,
		Description: description,
	}, nil
}

func readAmount(reader *bufio.Reader) (float64, error) {
	fmt.Print("Введите сумму: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)

	amount, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("сумма должна быть числом")
	}

	return amount, nil
}


func readText(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}



func main() {
	
	reader:= bufio.NewReader(os.Stdin)
	expenses := []Expense{}

	for {
		amount, err := readAmount(reader)

		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
	
		category, err := readText(reader, "Введите категорию: ")
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
		description, err := readText(reader, "Введите описание: ")
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
	
		expense, err := createExpense(amount, category, description)
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}

		expenses = append(expenses, expense)

		fmt.Println("Расход добавлен:")

		answer, err := readText(reader, "Добавить ещё? (y/n): ")

		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		if answer != "y" {
			break
		}
	}


	fmt.Println()
	fmt.Println("Ваши расходы:")


	//_, - range может вернуть две вещи: индекс + значение
	//_ означает: это значение я намеренно игнорирую.
	for _, expense := range expenses {
		fmt.Printf("%.2f ₽ | %s | %s\n",
			expense.Amount,
			expense.Category,
			expense.Description,
		)
	}
	
}

