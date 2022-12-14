package entity

type Expense struct {
	Category        string
	AmountInKopecks uint64
	Date            int64
}

type User struct {
	ID           int64
	CurrencyCode *string
	MonthlyLimit *uint64
	UpdatedAt    int64
}

type Report struct {
	Category        string
	AmountInKopecks uint64
}
