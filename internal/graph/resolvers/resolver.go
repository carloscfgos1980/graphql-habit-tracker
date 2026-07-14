package resolvers

import (
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	UserRepo *repository.UserRepository
}
