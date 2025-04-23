package repository

import (
	"context"
	"database/sql"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/models"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, tx *sql.Tx, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        INSERT INTO users (username, email, password, created_at, created_by, role_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;
    `

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRowContext(ctx, query,
			user.Username,
			user.Email,
			user.Password,
			user.CreatedAt,
			user.CreatedBy,
			user.RoleID,
		)
	} else {
		row = r.db.QueryRowContext(ctx, query,
			user.Username,
			user.Email,
			user.Password,
			user.CreatedAt,
			user.CreatedBy,
			user.RoleID,
		)
	}

	err := row.Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        SELECT u.id, 
               u.username, 
               u.email, 
               u.password, 
               u.is_active, 
               u.created_at, 
               u.created_by, 
               u.updated_at, 
               u.updated_by,
               r.id,
               r.name,
               r.level,
               r.description,
               r.is_active,
               r.created_at,
               r.created_by,
               r.updated_at,
               r.updated_by
        FROM users u
        JOIN roles r ON r.id = u.role_id
        WHERE email = $1;
    `

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
		&user.Role.IsActive,
		&user.Role.CreatedAt,
		&user.Role.CreatedBy,
		&user.Role.UpdatedAt,
		&user.Role.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, tx *sql.Tx, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        UPDATE users 
        SET username = $1, email = $2, password = $3, role_id = $4, is_active = $5, updated_at = $6, updated_by = $7
        WHERE id = $8
        RETURNING updated_at, updated_by;
    `

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role.ID,
			&user.IsActive,
			&user.UpdatedAt,
			&user.UpdatedBy,
			&user.ID,
		)
	} else {
		_, err = r.db.ExecContext(ctx, query,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.IsActive,
			&user.UpdatedAt,
			&user.UpdatedBy,
			&user.ID,
		)
	}
	if err != nil {
		return err
	}

	return nil
}
