package repository

import (
	"uas/app/models"
	"uas/database"

	"github.com/google/uuid"
)

func Login(loginInput string) (models.User, error) {
    var user models.User

    query := `
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.full_name, 
            u.role_id, r.name as role_name, u.is_active, u.created_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.username = $1 OR u.email = $1
    `

    err := database.ConnectDB().QueryRow(query, loginInput).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.FullName,
        &user.RoleID,
        &user.RoleName,
        &user.IsActive,
        &user.CreatedAt,
    )

    return user, err
}

func CheckPermission(userID uuid.UUID, permissionName string) (bool, error) {
    var exists bool

    query := `
        SELECT EXISTS (
            SELECT 1
            FROM users u
            JOIN roles r ON u.role_id = r.id
            JOIN role_permissions rp ON r.id = rp.role_id
            JOIN permissions p ON rp.permission_id = p.id
            WHERE u.id = $1 
              AND p.name = $2
        )
    `

    err := database.ConnectDB().QueryRow(query, userID, permissionName).Scan(&exists)
    if err != nil {
        return false, err
    }

    return exists, nil
}