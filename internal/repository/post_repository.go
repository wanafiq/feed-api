package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/models"
	"strings"
)

type PostRepository interface {
	Save(ctx context.Context, tx *sql.Tx, post *models.Post) error
	SavePostTag(ctx context.Context, tx *sql.Tx, postID string, tagName string) error
	SavePostUser(ctx context.Context, tx *sql.Tx, postID string, userID string) error
	FindAll(ctx context.Context, filter models.PostFilter) ([]*models.Post, int, error)
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

func (r *postRepository) FindAll(ctx context.Context, filter models.PostFilter) ([]*models.Post, int, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query, queryArgs := buildPostQuery(filter)
	fmt.Println("query", query)

	fmt.Println("queryArgs", queryArgs)

	countQuery, countArgs := buildPostCountQuery(filter)

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Content, &post.IsPublished,
			&post.PublishedAt, &post.CreatedAt, &post.CreatedBy,
			&post.UpdatedAt, &post.UpdatedBy, &post.AuthorID,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, &post)
	}

	return posts, total, rows.Err()
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

func buildPostQuery(filter models.PostFilter) (query string, args []any) {
	whereClauses, args := buildWhereClause(filter)
	orderBy := buildOrderByClause(filter.Sort)
	limitOffset := "LIMIT $%d OFFSET $%d"

	baseQuery := `
		SELECT
			id, title, slug, content, is_published, published_at,
			created_at, created_by, updated_at, updated_by, author_id
		FROM posts
	`

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Calculate the starting index for LIMIT and OFFSET placeholders
	limitOffsetStartIndex := len(args) + 1

	// Append LIMIT and OFFSET to the args slice
	args = append(args, filter.Limit, filter.Offset)

	finalQuery := baseQuery + fmt.Sprintf(" ORDER BY %s %s", orderBy, fmt.Sprintf(limitOffset, limitOffsetStartIndex, limitOffsetStartIndex+1))
	return finalQuery, args
}

func buildPostCountQuery(filter models.PostFilter) (query string, args []any) {
	whereClauses, args := buildWhereClause(filter)
	baseQuery := `SELECT COUNT(*) FROM posts`
	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	return baseQuery, args
}

func buildWhereClause(filter models.PostFilter) (whereClauses []string, args []any) {
	args = []any{}
	argID := 1
	whereClauses = []string{"1=1"} // Start with a condition that's always true

	// Search (title / content)
	if filter.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(title ILIKE '%%' || $%d || '%%' OR content ILIKE '%%' || $%d || '%%')", argID, argID+1))
		args = append(args, filter.Search, filter.Search)
		argID += 2
	}

	// Date range
	if filter.DateFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argID))
		args = append(args, *filter.DateFrom)
		argID++
	}
	if filter.DateTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argID))
		args = append(args, *filter.DateTo)
		argID++
	}

	// Tags (by name)
	if len(filter.Tags) > 0 {
		tagPlaceholders := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			tagPlaceholders[i] = fmt.Sprintf("$%d", argID)
			args = append(args, tag)
			argID++
		}
		whereClauses = append(whereClauses, fmt.Sprintf(`
			id IN (
				SELECT pt.post_id
				FROM post_tag pt
				JOIN tags t ON t.id = pt.tag_id
				WHERE t.name IN (%s)
				GROUP BY pt.post_id
				HAVING COUNT(DISTINCT t.name) = %d
			)
		`, strings.Join(tagPlaceholders, ", "), len(filter.Tags)))
	}

	return whereClauses, args
}

func buildOrderByClause(sort string) string {
	switch strings.ToLower(sort) {
	case "asc":
		return "created_at ASC"
	default:
		return "created_at DESC"
	}
}
