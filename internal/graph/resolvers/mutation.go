package resolvers

import (
	"context"
	"fmt"
	"os"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/graph/model"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/middleware"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/utils"
)

func (r *mutationResolver) Register(ctx context.Context, name string, email string, password string) (*model.AuthPayLoad, error) {
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

	return &model.AuthPayLoad{
		Token: token,
		User:  user,
	}, nil

}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayLoad, error) {
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

	return &model.AuthPayLoad{
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

// func (r *mutationResolver) DeleteUser(ctx context.Context) (bool, error) {
// 	userID, ok := middleware.GetUserID(ctx)

// 	if !ok {
// 		return false, fmt.Errorf("unauthorized")
// 	}

// 	deleted, err := r.UserRepo.DeleteUser(userID)
// 	if err != nil {
// 		return false, err
// 	}

// 	if !deleted {
// 		return false, fmt.Errorf("user not found")
// 	}

// 	return true, nil
// }
