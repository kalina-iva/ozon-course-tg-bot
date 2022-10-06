package entity

type Expense struct {
	Category string
	Amount   uint64
	Date     int64
}

type Report struct {
	Category string
	Amount   uint64
}
