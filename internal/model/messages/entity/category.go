package entity

type Category struct {
	Number int
	Name   string
}

type Expense struct {
	CategoryNumber int
	Amount         float64
	Date           int64
}
