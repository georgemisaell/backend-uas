package repository

import (
	"context"
	"database/sql"
	"time"
	"uas/app/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetByUsernameOrEmail(ctx context.Context, loginInput string) (models.User, error) // Tambahan buat Login
	CreateUser(ctx context.Context, tx *sql.Tx, user models.User) error // CreateUser biasanya butuh Transaction (Tx)
	UpdateUser(ctx context.Context, id uuid.UUID, user models.UpdateUser) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.FullName, &user.RoleID, &user.RoleName, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	var user models.User

	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.RoleName, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}

func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, loginInput string) (models.User, error) {
	var user models.User
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1 OR u.email = $1
	`
	err := r.db.QueryRowContext(ctx, query, loginInput).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.RoleName, 
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (r *userRepository) CreateUser(ctx context.Context, tx *sql.Tx, user models.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query,
			user.ID, user.Username, user.Email, user.PasswordHash,
			user.FullName, user.RoleID, user.IsActive, user.CreatedAt, user.UpdatedAt,
		)
	} else {
		_, err = r.db.ExecContext(ctx, query,
			user.ID, user.Username, user.Email, user.PasswordHash,
			user.FullName, user.RoleID, user.IsActive, user.CreatedAt, user.UpdatedAt,
		)
	}

	return err
}

func (r *userRepository) UpdateUser(ctx context.Context, id uuid.UUID, user models.UpdateUser) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, full_name = $3, role_id = $4, is_active = $5, updated_at = $6 
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.FullName, user.RoleID, user.IsActive, time.Now(), id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	query := `UPDATE users SET role_id = $1, updated_at = $2 WHERE id = $3`
	result, err := r.db.ExecContext(ctx, query, roleID, time.Now(), userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}