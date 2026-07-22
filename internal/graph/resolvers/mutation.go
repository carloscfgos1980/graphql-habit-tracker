package resolvers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/middleware"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/utils"
)

func (r *mutationResolver) Register(ctx context.Context, name string, email string, password string) (*models.AuthPayload, error) {
	err := utils.ValidateName(name)
	if err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}

	err = utils.ValidateEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	err = utils.ValidatePasswordStrength(password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := r.UserRepo.CreateUser(name, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthPayload{
		Token: token,
		User:  user,
	}, nil

}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*models.AuthPayload, error) {
	// Step 1: Fetch user by email
	user, err := r.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Step 2: Verify the password
	err = utils.ComparePassword(user.Password, password)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Step 3: Generate a JWT Token for the Authenticated User
	token, err := utils.GenerateJWT(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthPayload{
		Token: token,
		User:  user,
	}, nil

}

func (r *mutationResolver) UpdateUser(ctx context.Context, name *string, email *string, password *string) (*models.User, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	//fetch the user from the database to ensure they exist
	dbUser, err := r.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	newName := &dbUser.Username
	newEmail := &dbUser.Email
	newPassword := &dbUser.Password
	if name != nil {
		err := utils.ValidateName(*name)
		if err != nil {
			return nil, fmt.Errorf("invalid name: %w", err)
		}
		newName = name
	}

	if email != nil {
		if err := utils.ValidateEmail(*email); err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
		newEmail = email
	}

	if password != nil {
		if err := utils.ValidatePasswordStrength(*password); err != nil {
			return nil, fmt.Errorf("invalid password: %w", err)
		}

		hash, err := utils.HashPassword(*password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		newPassword = &hash
	}

	updatedUser, err := r.UserRepo.UpdateUser(userID, newName, newEmail, newPassword)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (bool, error) {
	userID, ok := middleware.GetUserID(ctx)

	if !ok {
		return false, fmt.Errorf("unauthorized")
	}

	deleted, err := r.UserRepo.DeleteUser(userID)
	if err != nil {
		return false, err
	}

	if !deleted {
		return false, fmt.Errorf("user not found")
	}

	return true, nil
}

// CreateHabit is the resolver for the createHabit field.
func (r *mutationResolver) CreateHabit(ctx context.Context, name string, description string) (*models.Habit, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	if err := utils.ValidateName(name); err != nil {
		return nil, fmt.Errorf("invalid habit name: %w", err)
	}
	if err := utils.ValidateDescription(description); err != nil {
		return nil, fmt.Errorf("invalid habit description: %w", err)
	}

	habit, err := r.HabitRepo.CreateHabit(userID, name, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	return habit, nil
}

// UpdateHabit is the resolver for the updateHabit field.
func (r *mutationResolver) UpdateHabit(ctx context.Context, id string, name *string, description *string) (*models.Habit, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	if name != nil {
		if err := utils.ValidateName(*name); err != nil {
			return nil, fmt.Errorf("invalid habit name: %w", err)
		}
	}
	if description != nil {
		if err := utils.ValidateDescription(*description); err != nil {
			return nil, fmt.Errorf("invalid habit description: %w", err)
		}
	}
	habit, err := r.HabitRepo.UpdateHabit(id, userID, name, description)
	if err != nil {
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	return habit, nil
}

// DeleteHabit is the resolver for the deleteHabit field.
func (r *mutationResolver) DeleteHabit(ctx context.Context, id string) (bool, error) {
	userID, ok := middleware.GetUserID(ctx)

	if !ok {
		return false, fmt.Errorf("unauthorized")
	}

	deleted, err := r.HabitRepo.DeleteHabit(id, userID)
	if err != nil {
		return false, err
	}

	if !deleted {
		return false, fmt.Errorf("habit not found")
	}

	return true, nil
}

// CheckInHabit is the resolver for the checkInHabit field.
func (r *mutationResolver) CheckInHabit(ctx context.Context, habitID string, date *string) (*models.HabitLog, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	var completedDate time.Time
	if date == nil {
		completedDate = time.Now().UTC()
	} else {
		var err error
		completedDate, err = time.Parse(time.RFC3339, *date)
		if err != nil {
			completedDate, err = time.Parse("2006-01-02", *date)
			if err != nil {
				return nil, fmt.Errorf("invalid date format: use RFC3339 or YYYY-MM-DD")
			}
			completedDate = completedDate.UTC()
		}
	}
	_, err := r.HabitRepo.GetHabitWithUserCheck(habitID, userID)
	if err != nil {
		return nil, err
	}
	isDuplicate, err := r.HabitLogRepo.CheckDuplicateLog(habitID, completedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate habit log: %w", err)
	}
	if isDuplicate {
		return nil, fmt.Errorf("duplicate habit log")
	}
	habitLog, err := r.HabitLogRepo.CreateHabitLog(habitID, completedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit log: %w", err)
	}
	return habitLog, nil
}

// DeleteHabitLog is the resolver for the deleteHabitLog field.
func (r *mutationResolver) DeleteHabitLog(ctx context.Context, id string) (bool, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return false, fmt.Errorf("unauthorized")
	}
	// fetch the habit log to ensure it exists and belongs to the user
	habitLog, err := r.HabitLogRepo.GetHabitLogByID(id)
	if err != nil {
		return false, fmt.Errorf("failed to fetch habit log: %w", err)
	}
	if habitLog == nil {
		return false, fmt.Errorf("habit log not found")
	}
	// check if the habit log belongs to a habit owned by the user
	habit, err := r.HabitRepo.GetHabitWithUserCheck(habitLog.HabitID, userID)
	if err != nil {
		return false, fmt.Errorf("unauthorized to delete this habit log: %w", err)
	}
	if habit == nil {
		return false, fmt.Errorf("unauthorized to delete this habit log")
	}
	deleted, err := r.HabitLogRepo.DeleteHabitLog(id)
	if err != nil {
		return false, fmt.Errorf("failed to delete habit log: %w", err)
	}
	if !deleted {
		return false, fmt.Errorf("habit log not found")
	}
	return true, nil
}
