package pgxhelper

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)


func (dh *DatabaseHelper) MakeTx(ctx context.Context) (pgx.Tx, func(err error) error, error) {
	tx, err := dh.Pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	cor := func(err error) error {
		if err != nil {
			if terr := tx.Rollback(ctx); terr != nil {
				err = errors.Wrap(terr, err.Error())
			}
			return err
		}
		return tx.Commit(ctx)
	}
	return tx, cor, nil

}