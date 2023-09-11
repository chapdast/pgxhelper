package pgxhelper

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (dh *DatabaseHelper) Exist(ctx context.Context, table string, condition string, values ...any) (bool, error) {
	exist := false
	var scanner Scanner = func(row pgx.Row) error {
		if err := row.Scan(&exist); err != nil {
			return err
		}
		return nil
	}

	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE %s )", table, condition)

	if err := dh.Query(ctx, scanner, query, values...); err != nil {
		return false, err
	}
	return exist, nil
}