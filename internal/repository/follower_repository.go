package repository

import (
	"context"
	"database/sql"
	"github.com/wanafiq/feed-api/internal/constants"
)

type FollowerRepository interface {
	Save(ctx context.Context, tx *sql.Tx, followerID string, followeeID string) error
	Delete(ctx context.Context, tx *sql.Tx, followerID string, followeeID string) error
}

type followerRepository struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) FollowerRepository {
	return &followerRepository{db: db}
}

func (r *followerRepository) Save(ctx context.Context, tx *sql.Tx, followerID string, followeeID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        INSERT INTO followers (follower_id, followee_id)
        VALUES ($1, $2);
    `

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, followerID, followeeID)
	} else {
		_, err = r.db.ExecContext(ctx, query, followerID, followeeID)
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *followerRepository) Delete(ctx context.Context, tx *sql.Tx, followerID string, followeeID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        DELETE FROM followers
        WHERE follower_id = $1 AND followee_id = $2;
    `

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, followerID, followeeID)
	} else {
		_, err = r.db.ExecContext(ctx, query, followerID, followeeID)
	}
	if err != nil {
		return err
	}

	return nil
}
