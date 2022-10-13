package entity

type Expense struct {
	Category        string
	AmountInKopecks uint64
	Date            int64
}

type Report struct {
	Category        string
	AmountInKopecks uint64
	Currency        string
}
