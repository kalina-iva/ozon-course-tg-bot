package entity

type Expense struct {
	Category        string
	AmountInKopecks uint64
	Date            int64
}

type User struct {
	UserID       int64
	CurrencyCode *string
	MonthlyLimit *uint64
	UpdatedAt    int64
}
