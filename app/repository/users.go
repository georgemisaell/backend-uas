package repository

import (
	"database/sql"
	"time"
	"uas/app/models"

	"github.com/google/uuid"
)

func GetAllUsers(db *sql.DB) ([]models.User, error){
	var users []models.User

	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
	`

	rows, err := db.Query(query)
	if err != nil{
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.RoleID,
			&user.RoleName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, err
}

func GetUserByID(db *sql.DB, id uuid.UUID) (models.User, error) {
	var user models.User

	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func CreateUser(tx *sql.Tx, user models.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := tx.Exec(query,
		user.ID,          
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func UpdateUser(db *sql.DB, id uuid.UUID, user models.UpdateUser) error {
	query := `
		UPDATE users 
		SET 
			username = $1, 
			email = $2, 
			full_name = $3, 
			role_id = $4, 
			is_active = $5, 
			updated_at = $6 
		WHERE id = $7
	`

	result, err := db.Exec(query,
		user.Username,
		user.Email,
		user.FullName,
		user.RoleID,
		user.IsActive,
		time.Now(),
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows 
	}

	return nil
}

func DeleteUser(db *sql.DB, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}