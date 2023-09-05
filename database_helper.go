package pgxhelper

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type DatabaseHelper struct {
	Pool *pgxpool.Pool
}

type Scanner = func(row pgx.Row) error

const (
	OPR_EQUAL = "="
	OPR_IN    = "IN"
	OPR_ANY   = "ANY"
	OPR_GTE   = ">="
)
const (
	JOINER_AND Joiner = " AND "
	JOINER_OR  Joiner = " OR "
)

type Joiner = string

type ConditionGroup struct {
	Join       Joiner
	Conditions []*Condition
	Group      *ConditionGroup
}

type Condition struct {
	ColumnName string
	Operator   string
	Value      any
}

var (
	ErrDuplicateKey = fmt.Errorf("DUPLICATE KEY")
	ErrNotFound = fmt.Errorf("NOT FOUND")
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

func (dh *DatabaseHelper) Query(ctx context.Context, scanner Scanner, query string, values ...any) error {
	tx, cor, err := dh.MakeTx(ctx)
	if err != nil {
		return err
	}
	rows, err := tx.Query(ctx, query, values...)
	if err != nil {
		return cor(err)
	}

	if scanner != nil {
		for rows.Next() {
			if err := scanner(rows); err != nil {
				return cor(err)
			}
		}
	}
	return cor(err)
}

func (dh *DatabaseHelper) Exist(ctx context.Context, table string, condition string, values ...any) (bool, error) {
	exist := false
	var scanner Scanner = func(row pgx.Row) error {
		if err := row.Scan(&exist); err != nil {
			return err
		}
		return nil
	}

	query := fmt.Sprintf("EXISTS (SELECT 1 FROM %s WHERE %s )", table, condition)

	if err := dh.Query(ctx, scanner, query, values...); err != nil {
		return false, err
	}
	return exist, nil
}


func (dh *DatabaseHelper) Exec(ctx context.Context, query string, values ...any) error {
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
