package memory

import (
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

func (m *Expense) New(userID int64, category string, amount uint64, date int64) {
	m.expenses[userID] = append(m.expenses[userID], &entity.Expense{
		Category:        category,
		AmountInKopecks: amount,
		Date:            date,
	})
}

func (m *Expense) Report(userID int64, period int64) []*entity.Report {
	reportMap := make(map[string]uint64)
	expenses := m.expenses[userID]
	for _, expense := range expenses {
		if expense.Date >= period {
			reportMap[expense.Category] += expense.AmountInKopecks
		}
	}
	report := make([]*entity.Report, 0, len(reportMap))
	for category, amount := range reportMap {
		report = append(report, &entity.Report{
			Category:        category,
			AmountInKopecks: amount,
		})
	}
	return report
}
