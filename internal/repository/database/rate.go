package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"log"
	"time"
)

type RateDB struct {
	conn *pgx.Conn
}

func NewRateDb(conn *pgx.Conn) *RateDB {
	return &RateDB{
		conn: conn,
	}
}

func (r *RateDB) GetRate(code string) (rate float64, err error) {
	rows, err := r.conn.Query(
		context.Background(),
		"select rate from exchange_rates where currency_code = $1 order by created_at desc limit 1",
		code,
	)
	if err != nil {
		log.Println("QueryRow failed:", err)
		return
	}
	defer rows.Close()

	has := rows.Next()
	if !has {
		err = errors.New("exchange rate not found by currency code")
		return
	}

	err = rows.Scan(&rate)

	return
}

func (r *RateDB) SaveRate(code string, rate float64) {
	_, err := r.conn.Exec(
		context.Background(),
		"insert into exchange_rates (currency_code, rate, created_at) VALUES ($1, $2, $3)",
		code,
		rate,
		time.Now(),
	)
	if err != nil {
		log.Println("cannot save rate:", err)
	}
}
