package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type RateDB struct {
	conn *pgx.Conn
}

func NewRateDb(conn *pgx.Conn) *RateDB {
	return &RateDB{
		conn: conn,
	}
}

func (r *RateDB) GetRate(ctx context.Context, code string) (float64, error) {
	row := r.conn.QueryRow(ctx, "select rate from exchange_rates where currency_code = $1 order by created_at desc limit 1", code)
	var rate float64
	err := row.Scan(&rate)
	return rate, err
}

func (r *RateDB) SaveRate(ctx context.Context, code string, rate float64) error {
	const sql = "insert into exchange_rates (currency_code, rate, created_at) VALUES ($1, $2, $3)"
	_, err := r.conn.Exec(ctx, sql, code, rate, time.Now())
	return err
}
