package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func withTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		e := fmt.Sprintf("failed to begin transaction: %s", err.Error())
		return errors.New(e)
	}

	err = fn(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			e := fmt.Sprintf("failed to rollback transaction: %s", err.Error())
			return errors.New(e)
		}
		return err
	}

	return tx.Commit()
}
