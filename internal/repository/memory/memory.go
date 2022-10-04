package memory

import (
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

type Expense struct {
	categories map[int64][]*entity.Category
	expenses   map[int64][]*entity.Expense
}

func New() *Expense {
	return &Expense{
		categories: make(map[int64][]*entity.Category),
		expenses:   make(map[int64][]*entity.Expense),
	}
}

func (m *Expense) NewCategory(userID int64, name string) *entity.Category {
	categories := m.categories[userID]

	cat := m.getCategoryByName(categories, name)
	if cat == nil {
		cat = &entity.Category{
			Number: calcNumber(categories),
			Name:   name,
		}

		m.categories[userID] = append(categories, cat)
	}
	return cat
}

func calcNumber(categories []*entity.Category) int {
	number := 0
	for _, cat := range categories {
		if cat.Number > number {
			number = cat.Number
		}
	}
	return number + 1
}

func (m *Expense) GetCategories(userID int64) []*entity.Category {
	categories := m.categories[userID]
	return categories
}

func (m *Expense) getCategoryByName(categories []*entity.Category, name string) *entity.Category {
	for _, cat := range categories {
		if cat.Name == name {
			return cat
		}
	}
	return nil
}

func (m *Expense) NewExpense(userID int64, category entity.Category, amount float64, date int64) {
	m.expenses[userID] = append(m.expenses[userID], &entity.Expense{
		Category: category,
		Amount:   amount,
		Date:     date,
	})
}

func (m *Expense) NewReport(userID int64, period int64) []*entity.Report {
	reportMap := make(map[entity.Category]float64)
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
