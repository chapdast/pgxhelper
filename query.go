package pgxhelper

import "context"

func (dh *DatabaseHelper) Query(ctx context.Context, scanner Scanner, query string, values ...any) error {
	/* tx, cor, err := dh.MakeTx(ctx)
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
	*/
			tx, cor, err := dh.MakeTx(ctx)
		if err != nil {
			return err
		}

		// Respect context cancelation before query
		if ctx.Err() != nil {
			return cor(ctx.Err())
		}

		rows, err := tx.Query(ctx, query, values...)
		if err != nil {
			return cor(err)
		}
		defer rows.Close()

		if scanner != nil {
			for rows.Next() {
				if ctx.Err() != nil {
					return cor(ctx.Err())
				}
				if err := scanner(rows); err != nil {
					return cor(err)
				}
			}
		}

		// Check rows.Err() to catch scan-level or driver-level issues
		if err := rows.Err(); err != nil {
			return cor(err)
		}

		// Final context check in case cancelation occurred during iteration
		if ctx.Err() != nil {
			return cor(ctx.Err())
		}

		return nil


}
