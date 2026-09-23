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
	Currency string
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

func createExpense(amount float64, category, description, currency string) (Expense, error) {
	if(amount<= 0) {
		return Expense{}, fmt.Errorf("сумма должна быть больше нуля")
	}


	if category == "" {
		return Expense{}, fmt.Errorf("категория не может быть пустой")
	}

	if description == "" {
		return Expense{}, fmt.Errorf("описание не может быть пустой")
	}


	if currency == "" {
		return Expense{}, fmt.Errorf("валюта не может быть пустой")
	}


	return Expense{
		Amount: amount,
		Category: category,
		Description: description,
		Currency:    currency,
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

//e — это конкретный Expense, с которым работает метод. как this в пхп
func (e Expense) IsValid() bool {
	return e.Amount > 0 && e.Category != ""
}

//* изменение объекта, указатель
func (e *Expense) UpdateDescription(description string) {
	e.Description = description
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

		currency, err := readText(reader, "Введите валюту: ")
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
	
		expense, err := createExpense(
			amount, 
			category, 
			description,
			currency,
		)
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
		fmt.Println(expense.Format())
	}
	
}

