package models

import "time"

type HabitLog struct {
	ID            string    `json:"id" db:"id"`
	HabitID       string    `json:"habit_id" db:"habit_id"`
	CompletedDate time.Time `json:"completed_date" db:"completed_date"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
