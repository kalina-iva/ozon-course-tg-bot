package memory

import (
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

type Expense struct {
	expenses map[int64][]*entity.Expense
}

func New() *Expense {
	return &Expense{
		expenses: make(map[int64][]*entity.Expense),
	}
}

func (m *Expense) NewExpense(userID int64, category string, amount int64, date int64) {
	m.expenses[userID] = append(m.expenses[userID], &entity.Expense{
		Category: category,
		Amount:   amount,
		Date:     date,
	})
}

func (m *Expense) NewReport(userID int64, period int64) []*entity.Report {
	reportMap := make(map[string]int64)
	expenses := m.expenses[userID]
	for _, expense := range expenses {
		if expense.Date >= period {
			reportMap[expense.Category] += expense.Amount
		}
	}
	report := make([]*entity.Report, 0, len(reportMap))
	for category, amount := range reportMap {
		report = append(report, &entity.Report{
			Category: category,
			Amount:   amount,
		})
	}
	return report
}
