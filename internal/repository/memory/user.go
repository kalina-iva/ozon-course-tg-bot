package memory

type User struct {
	currency map[int64]string
}

func NewUser() *User {
	return &User{
		currency: make(map[int64]string),
	}
}

func (u *User) SetCurrency(userID int64, currency string) error {
	u.currency[userID] = currency
	return nil
}

func (u *User) GetCurrency(userID int64) *string {
	currency, has := u.currency[userID]
	if has {
		return &currency
	}
	return nil
}
