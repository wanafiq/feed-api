package repository

import (
	"context"
	"database/sql"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/models"
)

type PostRepository interface {
	Save(ctx context.Context, tx *sql.Tx, post *models.Post) error
	SavePostTag(ctx context.Context, tx *sql.Tx, postID string, tagName string) error
	SavePostUser(ctx context.Context, tx *sql.Tx, postID string, userID string) error
	FindAll(ctx context.Context) ([]*models.Post, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*models.Post, error)
	FindByID(ctx context.Context, postID string) (*models.Post, error)
	Update(ctx context.Context, tx *sql.Tx, post *models.Post) error
	Delete(ctx context.Context, tx *sql.Tx, postID string) error
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Save(ctx context.Context, tx *sql.Tx, post *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        INSERT INTO posts (title, slug, content, is_published, published_at, created_at, created_by, author_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id;
    `

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRowContext(ctx, query,
			post.Title,
			post.Slug,
			post.Content,
			post.IsPublished,
			post.PublishedAt,
			post.CreatedAt,
			post.CreatedBy,
			post.AuthorID,
		)
	} else {
		row = r.db.QueryRowContext(ctx, query,
			post.Title,
			post.Slug,
			post.Content,
			post.IsPublished,
			post.PublishedAt,
			post.CreatedAt,
			post.CreatedBy,
			post.AuthorID,
		)
	}

	err := row.Scan(&post.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *postRepository) SavePostTag(ctx context.Context, tx *sql.Tx, postID string, tagID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        INSERT INTO post_tag (post_id, tag_id)
        VALUES ($1, $2);
    	`

	if tx != nil {
		_, err := tx.ExecContext(ctx, query, postID, tagID)
		if err != nil {
			return err
		}
	} else {
		_, err := r.db.ExecContext(ctx, query, postID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *postRepository) SavePostUser(ctx context.Context, tx *sql.Tx, postID string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        INSERT INTO post_user (post_id, user_id)
        VALUES ($1, $2);
    	`

	if tx != nil {
		_, err := tx.ExecContext(ctx, query, postID, userID)
		if err != nil {
			return err
		}
	} else {
		_, err := r.db.ExecContext(ctx, query, postID, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *postRepository) FindAll(ctx context.Context) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		SELECT 
			id, title, slug, content, is_published, published_at,
			created_at, created_by, updated_at, updated_by,
			author_id
		FROM posts
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Content,
			&post.IsPublished,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.CreatedBy,
			&post.UpdatedAt,
			&post.UpdatedBy,
			&post.AuthorID,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) FindAllByUserID(ctx context.Context, userID string) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	return nil, nil
}

func (r *postRepository) FindByID(ctx context.Context, postID string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	return nil, nil
}

func (r *postRepository) Update(ctx context.Context, tx *sql.Tx, post *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	return nil
}

func (r *postRepository) Delete(ctx context.Context, tx *sql.Tx, postID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	return nil
}
