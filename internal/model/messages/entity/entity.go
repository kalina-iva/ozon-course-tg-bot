package entity

type Category struct {
	Number int
	Name   string
}

type Expense struct {
	Category Category
	Amount   int64
	Date     int64
}

type Report struct {
	Category Category
	Amount   int64
}
