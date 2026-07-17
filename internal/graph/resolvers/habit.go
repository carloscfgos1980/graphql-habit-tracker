package resolvers

import (
	"context"
	"fmt"
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

// User is the resolver for the user field.
func (r *habitResolver) User(ctx context.Context, obj *models.Habit) (*models.User, error) {
	if obj == nil {
		return nil, fmt.Errorf("habit is required")
	}

	user, err := r.UserRepo.GetUserByID(obj.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve habit user: %w", err)
	}

	return user, nil
}
