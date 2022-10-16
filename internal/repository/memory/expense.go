package memory

import (
	"time"

	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

type Expense struct {
	expenses map[int64][]*entity.Expense
}

func NewExpense() *Expense {
	return &Expense{
		expenses: make(map[int64][]*entity.Expense),
	}
}

func (e *Expense) New(userID int64, category string, amount uint64, date time.Time) {
	e.expenses[userID] = append(e.expenses[userID], &entity.Expense{
		Category:        category,
		AmountInKopecks: amount,
		Date:            date.Unix(),
	})
}

func (e *Expense) GetExpenses(userID int64, period time.Time) []*entity.Expense {
	expenses := e.expenses[userID]
	periodUnix := period.Unix()
	var filtered []*entity.Expense
	for _, expense := range expenses {
		if expense.Date >= periodUnix {
			filtered = append(filtered, expense)
		}
	}
	return filtered
}
