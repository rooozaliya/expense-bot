package repository

import (
	"database/sql"
	"expense-bot/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) EnsureUser(chatID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`,
		chatID,
	)
	return err
}

func (r *ExpenseRepository) Add(chatID int64, e models.Expense) error {
	if err := r.EnsureUser(chatID); err != nil {
		return err
	}

	_, err := r.db.Exec(
		`INSERT INTO expenses (user_id, amount, category, description, currency, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		chatID, e.Amount, e.Category, e.Description, e.Currency, e.CreatedAt,
	)
	return err
}

func (r *ExpenseRepository) GetAll(chatID int64) ([]models.Expense, error) {
	rows, err := r.db.Query(
		`SELECT amount, category, description, currency, created_at
		 FROM expenses WHERE user_id = $1 ORDER BY created_at`,
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.Amount, &e.Category, &e.Description, &e.Currency, &e.CreatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}