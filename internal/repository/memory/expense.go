package memory

import (
	"time"

	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

type ExpenseM struct {
	expenses map[int64][]*entity.Expense
}

func NewExpenseM() *ExpenseM {
	return &ExpenseM{
		expenses: make(map[int64][]*entity.Expense),
	}
}

func (e *ExpenseM) New(userID int64, category string, amount uint64, date time.Time) error {
	e.expenses[userID] = append(e.expenses[userID], &entity.Expense{
		Category:        category,
		AmountInKopecks: amount,
		Date:            date.Unix(),
	})
	return nil
}

func (e *ExpenseM) Report(userID int64, period time.Time) []*entity.Report {
	reportMap := make(map[string]uint64)
	expenses := e.expenses[userID]
	periodUnix := period.Unix()
	for _, expense := range expenses {
		if expense.Date >= periodUnix {
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

func (e *ExpenseM) GetAmountByPeriod(userID int64, period time.Time) (sum uint64, err error) {
	expenses := e.expenses[userID]
	periodUnix := period.Unix()
	for _, expense := range expenses {
		if expense.Date >= periodUnix {
			sum += expense.AmountInKopecks
		}
	}
	return
}
