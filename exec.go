package pgxhelper

import (
	"context"
	"strings"
)

func (dh *DatabaseHelper) Exec(ctx context.Context, query string, values ...any) error {
	if tx, ok := dh.getTxOf(ctx); ok {
		result, err := tx.Exec(ctx, query, values...)
		if err != nil {
			if strings.HasPrefix(err.Error(), `ERROR: duplicate key value violates unique constraint`) {
				return ErrDuplicateKey
			}
			return err
		}
		if result.RowsAffected() == 0 {
			return ErrNotFound
		}
		return nil
	} else {
		tx, cor, err := dh.MakeTx(ctx)
		if err != nil {
			return err
		}
		result, err := tx.Exec(ctx, query, values...)
		if err != nil {
			if strings.HasPrefix(err.Error(), `ERROR: duplicate key value violates unique constraint`) {
				_ = cor(err)
				return ErrDuplicateKey
			}
			return cor(err)
		}
		if result.RowsAffected() == 0 {
			return ErrNotFound
		}
		return cor(err)
	}
}
