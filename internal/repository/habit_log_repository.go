package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
)

type HabitLogRepository struct {
	DB *sql.DB
}

func NewHabitLogRepository(db *sql.DB) *HabitLogRepository {
	return &HabitLogRepository{DB: db}
}

func (r *HabitLogRepository) CheckDuplicateLog(habitID string, completedDate time.Time) (bool, error) {
	// SQLITE's Date Column => 2026-06-06T10:30:00Z
	// stored date          => 2026-06-06
	// 2026-06-06T10:30:00Z ->>>>> 2026-06-06
	dateStr := completedDate.Format("2006-01-02")

	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM habit_logs WHERE habit_id = ? AND completed_date = ?", habitID, dateStr).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check duplicate log: %w", err)
	}

	return count > 0, nil
}

func (r *HabitLogRepository) CreateHabitLog(habitID string, completedDate time.Time) (*models.HabitLog, error) {
	now := time.Now()

	// 2026-06-06T10:30:00Z
	// 2026-06-06
	truncated := completedDate.UTC().Truncate(24 * time.Hour)
	dateStr := truncated.Format("2006-01-02")

	var id string
	err := r.DB.QueryRow("SELECT lower(hex(randomblob(16)))").Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate log ID: %w", err)
	}

	_, err = r.DB.Exec("INSERT INTO habit_logs (id, habit_id, completed_date, created_at) VALUES (?, ?, ?, ?)", id, habitID, dateStr, now)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("you have already checked in for this habit today or on this date")
		}
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	var log models.HabitLog
	err = r.DB.QueryRow("SELECT id, habit_id, completed_date, created_at FROM habit_logs WHERE id = ?", id).Scan(
		&log.ID,
		&log.HabitID,
		&log.CompletedDate,
		&log.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created log: %w", err)
	}

	return &log, nil
}

func (r *HabitLogRepository) GetHabitLogsByHabitID(habitID string) ([]*models.HabitLog, error) {
	rows, err := r.DB.Query("SELECT id, habit_id, completed_date, created_at FROM habit_logs WHERE habit_id = ? ORDER BY completed_date DESC", habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch habit logs: %w", err)
	}

	defer rows.Close()

	// nil -> json -> null
	// nil -> json -> []
	logs := make([]*models.HabitLog, 0)

	for rows.Next() {
		var log models.HabitLog
		err := rows.Scan(
			&log.ID,
			&log.HabitID,
			&log.CompletedDate,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit log: %w", err)
		}

		logs = append(logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habit logs: %w", err)
	}

	return logs, nil
}

func (r *HabitLogRepository) GetHabitLogByID(id string) (*models.HabitLog, error) {

	var log models.HabitLog
	err := r.DB.QueryRow("SELECT id, habit_id, completed_date, created_at FROM habit_logs WHERE id = ?", id).Scan(
		&log.ID,
		&log.HabitID,
		&log.CompletedDate,
		&log.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to fetch habit log: %w", err)
	}

	return &log, nil
}

func (r *HabitLogRepository) DeleteHabitLog(id string) (bool, error) {
	result, err := r.DB.Exec("DELETE FROM habit_logs WHERE id = ?", id)
	if err != nil {
		return false, fmt.Errorf("failed to delete habit log: %w", err)
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

func (r *HabitLogRepository) CountTotalCompletions(habitID string) (int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM habit_logs WHERE habit_id = ?", habitID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count completions: %w", err)
	}

	return count, nil
}
