package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
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

func (u *UserDB) GetUser(ctx context.Context, userID int64) (*entity.User, error) {
	var row pgx.Row
	tx := extractTx(ctx)
	if tx == nil {
		row = u.conn.QueryRow(ctx, "select currency_code, monthly_limit, updated_at from users where id=$1", userID)
	} else {
		row = tx.QueryRow(ctx, "select currency_code, monthly_limit, updated_at from users where id=$1", userID)
	}

	var currencyCode *string
	var monthlyLimit *uint64
	var updatedAt time.Time
	err := row.Scan(&currencyCode, &monthlyLimit, &updatedAt)

	switch err {
	case nil:
		return &entity.User{
			ID:           userID,
			CurrencyCode: currencyCode,
			MonthlyLimit: monthlyLimit,
			UpdatedAt:    updatedAt.Unix(),
		}, nil
	case pgx.ErrNoRows:
		return u.createUser(ctx, userID)
	default:
		return nil, errors.Wrap(err, "cannot scan user row")
	}
}

func (u *UserDB) createUser(ctx context.Context, userID int64) (*entity.User, error) {
	var err error
	timeNow := time.Now()
	tx := extractTx(ctx)
	if tx == nil {
		_, err = u.conn.Exec(ctx, "insert into users (id, updated_at) VALUES ($1, $2)", userID, timeNow)
	} else {
		_, err = tx.Exec(ctx, "insert into users (id, updated_at) VALUES ($1, $2)", userID, timeNow)
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot exec create user query")
	}
	return &entity.User{
		ID:           userID,
		CurrencyCode: nil,
		MonthlyLimit: nil,
		UpdatedAt:    timeNow.Unix(),
	}, nil
}

func (u *UserDB) SetCurrency(ctx context.Context, userID int64, currency string) error {
	return u.exec(ctx, "update users set currency_code = $1 where id = $2", currency, userID)
}

func (u *UserDB) SetLimit(ctx context.Context, userID int64, limit uint64) error {
	return u.exec(ctx, "update users set monthly_limit = $1 where id = $2", limit, userID)
}

func (u *UserDB) DelLimit(ctx context.Context, userID int64) error {
	return u.exec(ctx, "update users set monthly_limit = $1 where id = $2", nil, userID)
}

func (u *UserDB) exec(ctx context.Context, sql string, arguments ...any) error {
	var err error
	tx := extractTx(ctx)
	if tx == nil {
		_, err = u.conn.Exec(ctx, sql, arguments...)
	} else {
		_, err = tx.Exec(ctx, sql, arguments...)
	}
	return err
}
