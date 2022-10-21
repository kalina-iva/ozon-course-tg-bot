package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type TxManager struct {
	conn *pgx.Conn
}

func NewTxManager(conn *pgx.Conn) *TxManager {
	return &TxManager{
		conn: conn,
	}
}

type transactionKey struct{}

func (t *TxManager) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := t.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			log.Printf("rollback transaction: %v", errRollback)
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit(ctx); errCommit != nil {
		log.Printf("commit transaction: %v", errCommit)
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
