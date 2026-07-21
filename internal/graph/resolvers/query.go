package resolvers

import (
	"context"
	"fmt"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/middleware"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
)

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	userID, ok := middleware.GetUserID(ctx)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := r.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Habits is the resolver for the habits field.
func (r *queryResolver) Habits(ctx context.Context) ([]*models.Habit, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	habits, err := r.HabitRepo.GetHabitsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return habits, nil
}

// Habit is the resolver for the habit field.
func (r *queryResolver) Habit(ctx context.Context, id string) (*models.Habit, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	habit, err := r.HabitRepo.GetHabitWithUserCheck(id, userID)
	if err != nil {
		return nil, err
	}

	return habit, nil
}

// HabitLogs is the resolver for the habitLogs field.
func (r *queryResolver) HabitLogs(ctx context.Context, habitID string) ([]*models.HabitLog, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	_, err := r.HabitRepo.GetHabitWithUserCheck(habitID, userID)
	if err != nil {
		return nil, err
	}

	habitLogs, err := r.HabitLogRepo.GetHabitLogsByHabitID(habitID)
	if err != nil {
		return nil, err
	}

	return habitLogs, nil
}
