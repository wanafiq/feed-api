package repository

import (
	"context"
	"database/sql"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/models"
)

type RoleRepository interface {
	Save(ctx context.Context, tx *sql.Tx, role *models.Role) error
	FindByName(ctx context.Context, name string) (*models.Role, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Save(ctx context.Context, tx *sql.Tx, role *models.Role) error {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        INSERT INTO roles (name, level, description)
        VALUES ($1, $2, $3)
        RETURNING id;
    `

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRowContext(ctx, query,
			role.Name,
			role.Level,
			role.Description,
		)
	} else {
		row = r.db.QueryRowContext(ctx, query,
			role.Name,
			role.Level,
			role.Description,
		)
	}

	return row.Scan(&role.ID)
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	query := `
        SELECT id, name, level, description, is_active, created_at, created_by, updated_at, updated_by
        FROM roles
        WHERE name = $1;
    `

	role := &models.Role{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
		&role.IsActive,
		&role.CreatedAt,
		&role.CreatedBy,
		&role.UpdatedAt,
		&role.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}
