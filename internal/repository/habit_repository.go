package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
)

type HabitRepository struct {
	DB *sql.DB
}

func NewHabitRepository(db *sql.DB) *HabitRepository {
	return &HabitRepository{DB: db}
}

func (r *HabitRepository) CreateHabit(userID string, name string, description string) (*models.Habit, error) {
	now := time.Now()

	var id string
	err := r.DB.QueryRow("SELECT lower(hex(randomblob(16)))").Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate habit ID: %w", err)
	}

	_, err = r.DB.Exec("INSERT INTO habits (id, user_id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", id, userID, name, description, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create this habit: %w", err)
	}

	var habit models.Habit
	err = r.DB.QueryRow("SELECT id, user_id, name, description, created_at, updated_at FROM habits WHERE id = ?", id).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Description,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created habit: %w", err)
	}

	return &habit, nil
}

func (r *HabitRepository) GetHabitsByUserID(userID string) ([]*models.Habit, error) {
	// *sql.Rows => cursor
	rows, err := r.DB.Query("SELECT id, user_id, name, description, created_at, updated_at FROM habits WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch habits: %w", err)
	}

	defer rows.Close()

	// [] => OK
	// nil => NOT OK
	habits := make([]*models.Habit, 0)

	for rows.Next() {
		var habit models.Habit
		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Name,
			&habit.Description,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}

		habits = append(habits, &habit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}

func (r *HabitRepository) GetHabitWithUserCheck(habitID string, userID string) (*models.Habit, error) {
	var habit models.Habit

	err := r.DB.QueryRow("SELECT id, user_id, name, description, created_at, updated_at FROM habits WHERE id = ? AND user_id = ?", habitID, userID).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Description,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("habit not found")
		}

		return nil, fmt.Errorf("failed to fetch habit: %w", err)
	}

	return &habit, nil
}

func (r *HabitRepository) UpdateHabit(habitID string, userID string, name *string, description *string) (*models.Habit, error) {
	_, err := r.GetHabitWithUserCheck(habitID, userID)
	if err != nil {
		return nil, err
	}

	// SQL Fragment
	var setClauses []string
	// Data
	var args []interface{}

	if name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *name)
	}

	if description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *description)
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields provided to update")
	}

	now := time.Now()
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	setClause := strings.Join(setClauses, ", ")

	query := fmt.Sprintf("UPDATE habits SET %s WHERE id = ?", setClause)
	args = append(args, habitID)

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	return r.GetHabitWithUserCheck(habitID, userID)
}

func (r *HabitRepository) DeleteHabit(habitID string, userID string) (bool, error) {
	_, err := r.GetHabitWithUserCheck(habitID, userID)
	if err != nil {
		return false, err
	}

	result, err := r.DB.Exec("DELETE FROM habits WHERE id = ?", habitID)
	if err != nil {
		return false, fmt.Errorf("failed to delete habit: %w", err)
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

func (r *HabitRepository) GetHabitByID(id string) (*models.Habit, error) {
	var habit models.Habit

	err := r.DB.QueryRow("SELECT id, user_id, name, description, created_at, updated_at FROM habits WHERE id = ?", id).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Description,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to fetch habit: %w", err)
	}

	return &habit, nil
}
