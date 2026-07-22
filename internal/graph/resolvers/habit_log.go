package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/middleware"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
)

// CompletedDate is the resolver for the completedDate field.
func (r *habitLogResolver) CompletedDate(ctx context.Context, obj *models.HabitLog) (string, error) {
	return obj.CompletedDate.Format("2006-01-02"), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *habitLogResolver) CreatedAt(ctx context.Context, obj *models.HabitLog) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// Habit is the resolver for the habit field.
func (r *habitLogResolver) Habit(ctx context.Context, obj *models.HabitLog) (*models.Habit, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	habit, err := r.HabitRepo.GetHabitWithUserCheck(obj.HabitID, userID)
	if err != nil {
		return nil, err
	}

	return habit, nil
}
