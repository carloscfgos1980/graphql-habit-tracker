package resolvers

import (
	"context"
	"time"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
)

// CreatedAt is the resolver for the createdAt field.
func (r *habitResolver) CreatedAt(ctx context.Context, obj *models.Habit) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *habitResolver) UpdatedAt(ctx context.Context, obj *models.Habit) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil

}
