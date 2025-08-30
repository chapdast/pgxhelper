package pgxhelper

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type KEY string

const KeyTransactional KEY = "__TX-ACTION__"

type CRFunc = func(error) error

func makeCor(ctx context.Context, tx pgx.Tx) CRFunc {
	return func(err error) error {
		if err != nil {
			if terr := tx.Rollback(ctx); terr != nil {
				err = errors.Wrap(terr, err.Error())
			}
			return err
		}
		return tx.Commit(ctx)
	}
}
func (dh *DatabaseHelper) MakeTx(ctx context.Context) (pgx.Tx, CRFunc, error) {
	tx, err := dh.Pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	cor := makeCor(ctx, tx)
	return tx, cor, nil

}

var (
	ErrAlreadyTransactional = fmt.Errorf("alreay a transaction")
)

func (dh *DatabaseHelper) isTransactional(ctx context.Context) bool {
	_, ok := ctx.Value(KeyTransactional).(pgx.Tx)
	return ok
}
func (dh *DatabaseHelper) getTxOf(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(KeyTransactional).(pgx.Tx)
	return tx, ok

}
func (dh *DatabaseHelper) Transactional(ctx context.Context) (context.Context, CRFunc, error) {

	if !dh.isTransactional(ctx) {
		tx, err := dh.Pool.Begin(ctx)
		if err != nil {
			return nil, nil, err
		}
		ctx = context.WithValue(ctx, KeyTransactional, tx)
		return ctx, makeCor(ctx, tx), nil
	}
	return nil, nil, ErrAlreadyTransactional
}
