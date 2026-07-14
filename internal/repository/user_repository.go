package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
)

// connection pool -> a collection of db connections
type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(username string, email string, hashedPassword string) (*models.User, error) {
	now := time.Now()

	var id string
	err := r.DB.QueryRow("SELECT lower(hex(randomblob(16)))").Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to generate user ID: %w", err)
	}

	_, err = r.DB.Exec("INSERT INTO users (id, username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", id, username, email, hashedPassword, now, now)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("email already registered")
		}

		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	var user models.User
	err = r.DB.QueryRow("SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?", id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = ?", email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid email or password")
		}

		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil

}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?", id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(id string, username *string, email *string, password *string) (*models.User, error) {
	// SQL Fragment
	var setClauses []string
	// Data
	var args []interface{}

	// SQL Statement
	// UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?

	// UPDATE users SET name = ? WHERE id = ?
	// UPDATE users SET email = ? WHERE id = ?
	// UPDATE users SET password = ? WHERE id = ?
	// UPDATE users SET name = ?, email = ? WHERE id = ?

	// setClauses -> ["name = ?"]
	// args       -> ["Samantha"]

	// setClauses -> ["email = ?"]
	// args       -> ["Samantha-smith@test.com"]

	// setClauses -> ["password = ?"]
	// args       -> ["MyPass@13456"]

	// setClauses -> ["name = ?", "email = ?"]
	// args       -> ["Samantha", "Samantha-smith@test.com"]

	if username != nil {
		setClauses = append(setClauses, "username = ?")
		args = append(args, *username)
	}

	if email != nil {
		setClauses = append(setClauses, "email = ?")
		args = append(args, *email)
	}

	if password != nil {
		setClauses = append(setClauses, "password = ?")
		args = append(args, *password)
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields provided to update")
	}

	now := time.Now()
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	setClause := strings.Join(setClauses, ", ")

	// SQL Statement/Query
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", setClause)
	args = append(args, id)

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("email already in use")
		}

		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.GetUserByID(id)
}

func (r *UserRepository) DeleteUser(id string) (bool, error) {
	result, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to confirm deletion: %w", err)

	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
