package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/models"
	"strings"
	"time"
)

type PostRepository interface {
	Save(ctx context.Context, tx *sql.Tx, post *models.Post) error
	FindAll(ctx context.Context, filter models.PostFilter) ([]*models.Post, int, error)
	FindByID(ctx context.Context, postID string) (*models.Post, error)
	Update(ctx context.Context, tx *sql.Tx, post *models.Post) (*models.Post, error)
	Delete(ctx context.Context, tx *sql.Tx, postID string) error

	SavePostTag(ctx context.Context, tx *sql.Tx, postID string, tagName string) error
	DeletePostTag(ctx context.Context, tx *sql.Tx, postID string) error

	SavePostUser(ctx context.Context, tx *sql.Tx, postID string, userID string) error
	DeletePostUser(ctx context.Context, tx *sql.Tx, postID string) error
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{
		db: db,
	}
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

	query, countQuery, queryArgs, countArgs := buildPostQuery(filter)

	// total count query
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// data query (with pagination)
	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
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
			&post.Author.ID,
			&post.Author.Username,
			&post.Author.Email,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, &post)
	}

	return posts, total, rows.Err()
}

func (r *postRepository) FindByID(ctx context.Context, postID string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		SELECT
			p.id, p.title, p.slug, p.content, p.is_published, p.published_at,
			p.created_at, p.created_by, p.updated_at, p.updated_by, p.author_id,
			u.id as author_id, u.username, u.email as author_email,
			r.id as role_id, r.name as role_name, r.level as role_level, r.description as role_description, 
			r.is_active as role_is_active, r.created_at as role_created_at, r.created_by as role_created_by, 
			r.updated_at as role_updated_at, r.updated_by as role_updated_by
		FROM posts p
		JOIN users u ON p.author_id = u.id
		JOIN roles r ON u.role_id = r.id
		WHERE p.id = $1
	`

	row := r.db.QueryRowContext(ctx, query, postID)

	var post models.Post
	var author models.User
	var role models.Role
	err := row.Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content, &post.IsPublished,
		&post.PublishedAt, &post.CreatedAt, &post.CreatedBy, &post.UpdatedAt, &post.UpdatedBy, &post.AuthorID,
		&author.ID, &author.Username, &author.Email,
		&role.ID, &role.Name, &role.Level, &role.Description, &role.IsActive, &role.CreatedAt, &role.CreatedBy, &role.UpdatedAt, &role.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}
	post.Author = author
	post.Author.Role = role

	return &post, nil
}

func (r *postRepository) Update(ctx context.Context, tx *sql.Tx, post *models.Post) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		UPDATE posts
		SET
			author_id = $1,
			title = $2,
			slug = $3,
			content = $4,
			is_published = $5,
			published_at = $6,
			updated_at = $7,
			updated_by = $8
		WHERE id = $9
		RETURNING id, author_id, title, slug, content, is_published, published_at, created_at, created_by, updated_at, updated_by;
	`

	updatedPost := &models.Post{}
	var row *sql.Row

	if tx != nil {
		row = tx.QueryRowContext(
			ctx,
			query,
			post.AuthorID,
			post.Title,
			post.Slug,
			post.Content,
			post.IsPublished,
			post.PublishedAt,
			time.Now().UTC(),
			post.UpdatedBy,
			post.ID,
		)
	} else {
		row = r.db.QueryRowContext(
			ctx,
			query,
			post.AuthorID,
			post.Title,
			post.Slug,
			post.Content,
			post.IsPublished,
			post.PublishedAt,
			time.Now().UTC(),
			post.UpdatedBy,
			post.ID,
		)
	}

	err := row.Scan(
		&updatedPost.ID, &updatedPost.AuthorID, &updatedPost.Title, &updatedPost.Slug, &updatedPost.Content,
		&updatedPost.IsPublished, &updatedPost.PublishedAt, &updatedPost.CreatedAt, &updatedPost.CreatedBy,
		&updatedPost.UpdatedAt, &updatedPost.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

func (r *postRepository) Delete(ctx context.Context, tx *sql.Tx, postID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		DELETE FROM posts
		WHERE id = $1
	`

	var err error

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, postID)
	} else {
		_, err = r.db.ExecContext(ctx, query, postID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *postRepository) DeletePostTag(ctx context.Context, tx *sql.Tx, postID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		DELETE FROM post_tag
		WHERE post_id = $1
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, postID)
	} else {
		_, err = r.db.ExecContext(ctx, query, postID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *postRepository) DeletePostUser(ctx context.Context, tx *sql.Tx, postID string) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
		DELETE FROM post_user
		WHERE post_id = $1
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, postID)
	} else {
		_, err = r.db.ExecContext(ctx, query, postID)
	}

	if err != nil {
		return err
	}

	return nil
}

func buildPostQuery(filter models.PostFilter) (query string, countQuery string, queryArgs []any, countArgs []any) {
	var baseArgs []any
	argID := 1
	where := []string{"1=1"}

	// Search in title or content
	if filter.Search != "" {
		where = append(where,
			fmt.Sprintf("(p.title ILIKE '%%' || $%d || '%%' OR p.content ILIKE '%%' || $%d || '%%')", argID, argID),
		)
		baseArgs = append(baseArgs, filter.Search)
		argID++
	}

	// Date filters
	if filter.DateFrom != nil {
		where = append(where, fmt.Sprintf("p.created_at >= $%d", argID))
		baseArgs = append(baseArgs, *filter.DateFrom)
		argID++
	}
	if filter.DateTo != nil {
		where = append(where, fmt.Sprintf("p.created_at <= $%d", argID))
		baseArgs = append(baseArgs, *filter.DateTo)
		argID++
	}

	// Tag filtering
	if len(filter.Tags) > 0 {
		tagPlaceholders := []string{}
		for _, tag := range filter.Tags {
			tagPlaceholders = append(tagPlaceholders, fmt.Sprintf("$%d", argID))
			baseArgs = append(baseArgs, tag)
			argID++
		}
		where = append(where, fmt.Sprintf(`
			p.id IN (
				SELECT pt.post_id
				FROM post_tag pt
				JOIN tags t ON t.id = pt.tag_id
				WHERE t.name IN (%s)
				GROUP BY pt.post_id
				HAVING COUNT(DISTINCT t.name) = %d
			)
		`, strings.Join(tagPlaceholders, ", "), len(filter.Tags)))
	}

	// WHERE clause
	whereClause := strings.Join(where, " AND ")

	// Select with JOIN on users table
	selectFields := `
		SELECT 
			p.id, p.title, p.slug, p.content, p.is_published, p.published_at,
			p.created_at, p.created_by, p.updated_at, p.updated_by, p.author_id,
			u.id, u.username, u.email
		FROM posts p
		JOIN users u ON u.id = p.author_id
		WHERE ` + whereClause

	// Count query doesn't need user fields
	countQuery = `SELECT COUNT(*) FROM posts p WHERE ` + whereClause
	countArgs = append(countArgs, baseArgs...)

	// Sorting
	order := "p.created_at DESC"
	if strings.ToLower(filter.Sort) == "asc" {
		order = "p.created_at ASC"
	}

	// Pagination placeholders
	limitPlaceholder := fmt.Sprintf("$%d", argID)
	offsetPlaceholder := fmt.Sprintf("$%d", argID+1)
	queryArgs = append(queryArgs, baseArgs...)
	queryArgs = append(queryArgs, filter.Limit, filter.Offset)

	// Final query
	query = fmt.Sprintf(`%s ORDER BY %s LIMIT %s OFFSET %s`, selectFields, order, limitPlaceholder, offsetPlaceholder)

	return query, countQuery, queryArgs, countArgs
}
