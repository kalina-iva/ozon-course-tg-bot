package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type transactionKey struct{}

type TxManager struct {
	conn *pgx.Conn
}

func NewTxManager(conn *pgx.Conn) *TxManager {
	return &TxManager{
		conn: conn,
	}
}

func (t *TxManager) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := t.conn.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}

	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			log.Println("rollback transaction:", errRollback)
		}
		return err
	}
	if errCommit := tx.Commit(ctx); errCommit != nil {
		log.Println("commit transaction:", errCommit)
	}
	return nil
}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, transactionKey{}, tx)
}

func extractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(transactionKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}