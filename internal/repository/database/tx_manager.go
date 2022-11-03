package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
			zap.L().Error("rollback transaction failed", zap.Error(errRollback))
		}
		return errors.Wrap(err, "cannot exec tFunc")
	}
	if errCommit := tx.Commit(ctx); errCommit != nil {
		zap.L().Error("commit transaction failed", zap.Error(errCommit))
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			zap.L().Error("rollback transaction failed", zap.Error(errRollback))
		}
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
