package pgxhelper

import "context"

func (dh *DatabaseHelper) Query(ctx context.Context, scanner Scanner, query string, values ...any) error {
	if tx, ok := dh.getTxOf(ctx); ok {
		rows, err := tx.Query(ctx, query, values...)
		if err != nil {
			return err
		}
		if scanner != nil {
			for rows.Next() {
				if err := scanner(rows); err != nil {
					return err
				}
			}
		}
		return nil
	} else {
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
}
