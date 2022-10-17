package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

type ExpenseDB struct {
	conn *pgx.Conn
}

func NewExpenseDb(conn *pgx.Conn) *ExpenseDB {
	return &ExpenseDB{
		conn: conn,
	}
}

func (e *ExpenseDB) New(userID int64, category string, amount uint64, date time.Time) error {
	_, err := e.conn.Exec(
		context.Background(),
		"insert into expenses (user_id, category, amount, created_at) VALUES ($1, $2, $3, $4)",
		userID,
		category,
		amount,
		date,
	)
	return err
}

func (e *ExpenseDB) GetExpenses(userID int64, period time.Time) []*entity.Expense {
	rows, err := e.conn.Query(
		context.Background(),
		"select category, amount, created_at from expenses where user_id = $1 and created_at >= $2",
		userID,
		period,
	)
	if err != nil {
		log.Println("cannot exec get expenses query:", err)
		return nil
	}
	defer rows.Close()

	var expenses []*entity.Expense
	for rows.Next() {
		var category string
		var amount uint64
		var createdAt time.Time
		if err := rows.Scan(&category, &amount, &createdAt); err != nil {
			log.Println("cannot scan expense row:", err)
			break
		}

		expenses = append(expenses, &entity.Expense{
			Category:        category,
			AmountInKopecks: amount,
			Date:            createdAt.Unix(),
		})
	}
	return expenses
}
