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

func (m *Expense) NewCategory(userId int64, name string) *entity.Category {
	categories := m.categories[userId]

	cat := m.getCategoryByName(categories, name)
	if cat == nil {
		cat = &entity.Category{
			Number: calcNumber(categories),
			Name:   name,
		}

		m.categories[userId] = append(categories, cat)
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

func (m *Expense) GetCategories(userId int64) []*entity.Category {
	categories := m.categories[userId]
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

func (m *Expense) NewExpense(userId int64, categoryNumber int, amount float64, date int64) {
	m.expenses[userId] = append(m.expenses[userId], &entity.Expense{
		CategoryNumber: categoryNumber,
		Amount:         amount,
		Date:           date,
	})
}
