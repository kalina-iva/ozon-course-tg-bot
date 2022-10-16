package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

type UserDB struct {
	conn *pgx.Conn
}

func NewUserDb(conn *pgx.Conn) *UserDB {
	return &UserDB{
		conn: conn,
	}
}

func (u *UserDB) GetUser(userID int64) *entity.User {
	rows, err := u.conn.Query(
		context.Background(),
		"select currency_code, monthly_limit, updated_at from users where user_id=$1",
		userID,
	)
	if err != nil {
		log.Fatal("QueryRow failed:", err)
	}
	defer rows.Close()

	var user *entity.User
	for rows.Next() {
		var currencyCode *string
		var monthlyLimit *uint64
		var updatedAt time.Time
		if err := rows.Scan(&currencyCode, &monthlyLimit, &updatedAt); err != nil {
			log.Fatal("Row scan failed:", err)
		}

		user = &entity.User{
			UserID:       userID,
			CurrencyCode: currencyCode,
			MonthlyLimit: monthlyLimit,
			UpdatedAt:    updatedAt.Unix(),
		}
	}
	return user
}

func (u *UserDB) CreateUser(userID int64, currency *string, monthlyLimit *uint64) error {
	s, err := u.conn.Exec(
		context.Background(),
		"insert into users (user_id, currency_code, monthly_limit, updated_at) VALUES ($1, $2, $3, $4)",
		userID,
		currency,
		monthlyLimit,
		time.Now(),
	)
	fmt.Printf("Res: %v\n", s)
	return err
}

func (u *UserDB) SetCurrency(userID int64, currency string) error {
	user := u.GetUser(userID)
	var err error
	if user == nil {
		err = u.CreateUser(userID, &currency, nil)
	} else {
		_, err = u.conn.Exec(context.Background(), "update users set currency_code = $1 where user_id = $2", currency, userID)
	}

	return err
}

func (u *UserDB) GetCurrency(userID int64) *string {
	user := u.GetUser(userID)
	if user != nil {
		return user.CurrencyCode
	}
	return nil
}

func (u *UserDB) SetLimit(userID int64, limit uint64) error {
	user := u.GetUser(userID)
	var err error
	if user == nil {
		err = u.CreateUser(userID, nil, &limit)
	} else {
		_, err = u.conn.Exec(context.Background(), "update users set monthly_limit = $1 where user_id = $2", limit, userID)
	}

	return err
}

func (u *UserDB) GetLimit(userID int64) *uint64 {
	user := u.GetUser(userID)
	if user != nil {
		return user.MonthlyLimit
	}
	return nil
}

func (u *UserDB) DelLimit(userID int64) error {
	user := u.GetUser(userID)
	var err error
	if user != nil {
		_, err = u.conn.Exec(context.Background(), "update users set monthly_limit = $1 where user_id = $2", nil, userID)
	}
	return err
}
