package memory

type UserM struct {
	currency map[int64]string
}

func NewUserM() *UserM {
	return &UserM{
		currency: make(map[int64]string),
	}
}

func (u *UserM) SetCurrency(userID int64, currency string) error {
	u.currency[userID] = currency
	return nil
}

func (u *UserM) GetCurrency(userID int64) *string {
	currency, has := u.currency[userID]
	if has {
		return &currency
	}
	return nil
}

func (u *UserM) SetLimit(userID int64, limit uint64) error {
	return nil
}

func (u *UserM) DelLimit(userID int64) error {
	return nil
}

func (u *UserM) GetLimit(userID int64) *uint64 {
	return nil
}
