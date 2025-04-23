package repository

import (
	"context"
	"database/sql"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/models"
)

type TokenRepository interface {
	Save(ctx context.Context, tx *sql.Tx, role *models.Token) error
	FindByUserID(ctx context.Context, userID string) (*models.Token, error)
}

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Save(ctx context.Context, tx *sql.Tx, token *models.Token) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		INSERT INTO tokens (type, value, expired_at, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRowContext(ctx, query,
			token.Type,
			token.Value,
			token.ExpiredAt,
			token.UserID,
		)
	} else {
		row = r.db.QueryRowContext(ctx, query,
			token.Type,
			token.Value,
			token.ExpiredAt,
			token.UserID,
		)
	}

	err := row.Scan(&token.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *tokenRepository) FindByUserID(ctx context.Context, userID string) (*models.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		SELECT id, type, value, expired_at
		FROM tokens
		WHERE user_id = $1;
	`

	token := &models.Token{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&token.ID,
		&token.Type,
		&token.Value,
		&token.ExpiredAt,
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}
