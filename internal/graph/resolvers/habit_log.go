package resolvers

import (
	"context"
	"time"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
)

// CompletedDate is the resolver for the completedDate field.
func (r *habitLogResolver) CompletedDate(ctx context.Context, obj *models.HabitLog) (string, error) {
	return obj.CompletedDate.Format(time.RFC3339), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *habitLogResolver) CreatedAt(ctx context.Context, obj *models.HabitLog) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}
