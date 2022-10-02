package memory

import (
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

type Expense struct {
	categories map[int64][]*entity.Category
}

func New() *Expense {
	return &Expense{
		categories: make(map[int64][]*entity.Category),
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
