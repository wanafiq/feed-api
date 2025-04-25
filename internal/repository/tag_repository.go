package repository

import (
	"context"
	"database/sql"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/models"
)

type TagRepository interface {
	Save(ctx context.Context, tx *sql.Tx, tag *models.Tag) error
	FindAll(ctx context.Context) ([]*models.Tag, error)
	FindByName(ctx context.Context, tagName string) (*models.Tag, error)
	FindByID(ctx context.Context, tagID string) (*models.Tag, error)
	FindByPostID(ctx context.Context, postID string) ([]*models.Tag, error)
	Delete(ctx context.Context, tx *sql.Tx, tagID string) error
}

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Save(ctx context.Context, tx *sql.Tx, tag *models.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        INSERT INTO tags (name)
        VALUES ($1)
        RETURNING id;
    `
	var row *sql.Row

	if tx != nil {
		row = tx.QueryRowContext(ctx, query, tag.Name)
	} else {
		row = r.db.QueryRowContext(ctx, query, tag.Name)
	}

	err := row.Scan(&tag.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *tagRepository) FindAll(ctx context.Context) ([]*models.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	return nil, nil
}

func (r *tagRepository) FindByName(ctx context.Context, tagName string) (*models.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        SELECT id, name
        FROM tags
        WHERE name = $1;
    `

	tag := &models.Tag{}
	err := r.db.QueryRowContext(ctx, query, tagName).Scan(&tag.ID, &tag.Name)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *tagRepository) FindByID(ctx context.Context, tagID string) (*models.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	return nil, nil
}

func (r *tagRepository) FindByPostID(ctx context.Context, postID string) ([]*models.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		SELECT
			t.id, t.name
		FROM tags t
		JOIN post_tag pt ON t.id = pt.tag_id
		WHERE pt.post_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, rows.Err()
}

func (r *tagRepository) Delete(ctx context.Context, tx *sql.Tx, tagID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	return nil
}
