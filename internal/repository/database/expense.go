package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

type ExpenseDB struct {
	conn *pgx.Conn
}

func NewExpenseDb(conn *pgx.Conn) *ExpenseDB {
	return &ExpenseDB{
		conn: conn,
	}
}

func (e *ExpenseDB) New(ctx context.Context, userID int64, category string, amount uint64, date time.Time) error {
	const sql = "insert into expenses (user_id, category, amount, created_at) VALUES ($1, $2, $3, $4)"
	tx := extractTx(ctx)
	var err error
	if tx == nil {
		_, err = e.conn.Exec(ctx, sql, userID, category, amount, date)
	} else {
		_, err = tx.Exec(ctx, sql, userID, category, amount, date)
	}
	return errors.Wrap(err, "cannot save expense")
}

func (e *ExpenseDB) Report(ctx context.Context, userID int64, period time.Time) ([]*entity.Report, error) {
	const sql = "select category, sum(amount) as sum from expenses where user_id = $1 and created_at >= $2 group by category"
	var rows pgx.Rows
	var err error
	tx := extractTx(ctx)
	if tx == nil {
		rows, err = e.conn.Query(ctx, sql, userID, period)
	} else {
		rows, err = tx.Query(ctx, sql, userID, period)
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot exec get expenses query")
	}
	defer rows.Close()

	var reportData []*entity.Report
	for rows.Next() {
		var report entity.Report
		if err = rows.Scan(&report.Category, &report.AmountInKopecks); err != nil {
			return nil, errors.Wrap(err, "cannot scan report row")
		}
		reportData = append(reportData, &report)
	}
	return reportData, nil
}

func (e *ExpenseDB) GetAmountByPeriod(ctx context.Context, userID int64, period time.Time) (uint64, error) {
	const sql = "select sum(amount) as sum from expenses where user_id = $1 and created_at >= $2 group by user_id"
	tx := extractTx(ctx)
	var row pgx.Row
	if tx == nil {
		row = e.conn.QueryRow(ctx, sql, userID, period)
	} else {
		row = tx.QueryRow(ctx, sql, userID, period)
	}
	var sum uint64
	err := row.Scan(&sum)
	if err != nil {
		return 0, errors.Wrap(err, "cannot scan sum row")
	}
	return sum, nil
}
