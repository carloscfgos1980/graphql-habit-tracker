package models

import "time"

type Habit struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	CurrentStreak    int `json:"current_streak"`
	LongestStreak    int `json:"longest_streak"`
	TotalCompletions int `json:"total_completions"`
}
