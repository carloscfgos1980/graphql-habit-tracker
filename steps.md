# STEPS

[Go GraphQL API for Beginner](https://www.youtube.com/watch?v=rrY7tcDSGZ8)

go get github.com/gin-gonic/gin
go get github.com/mattn/go-sqlite3
go get github.com/joho/godotenv

go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen init
go run github.com/99designs/gqlgen generate

- Move graph into internal directory

go install github.com/air-verse/air@latest
export PATH="$PATH:$(go env GOPATH)/bin"
# Needed because `go install` puts `air` in GOPATH/bin, which is not always on zsh's PATH.


goose -dir ./migrations sqlite3 ./data/habit.db up
goose -dir ./migrations sqlite3 ./data/habit.db status
goose -dir ./migrations sqlite3 ./data/habit.db down

## migrations
- Create table for users
- Create table for habits with index in user_id
- Create table for habit logs with index of habit_id and completed_date
- Run CLI to create tables with goose
goose -dir ./migrations sqlite3 ./data/habit.db up

## schema.graphqls
* this is pretty important. Here is define the type of answers to the resquests, the mutations and the queries
- run de script to generate schema.resolver.go file:
go run github.com/99designs/gqlgen generate


## main
- Load environment variables from .env file
- get the database path from environment variable or use default
- Initialize the database
- Initialize repositories
- Initialize GraphQL server
- Initialize Gin router
- GraphQL Playground and GraphQL endpoint
- Apply authentication middleware to the /graphql endpoint
- get the port from environment variable
- Start the server and log any errors

## middleware
>internal/middleware/auth_middleware
- AuthMiddleware is a Gin middleware that checks for a valid JWT token in the Authorization header.
- GetUserID retrieves the user ID from the context. It returns the user ID and a boolean indicating whether the user ID was found.

## register

1. auxiliar functions:
- email
- password
- description
- jwt

2. Repository. 
CreateUser creates a new user in the database with the provided username, email, and hashed password.

3. Register resolves the register mutation, creating a new user and returning an authentication payload.

4. 	Initialize user repository
	userRepo := repository.NewUserRepository(db)

5. Add the repository to graphql server
	graphqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{
			UserRepo:     userRepo,
		},
	}))

## login
1. user_repository
GetUserByEmail retrieves a user from the database by their email address. It returns the user and an error if any.
2. Mutations
Login resolves the login mutation, authenticating a user and returning an authentication payload.

## update user
1. user_repository
- GetUserByID retrieves a user from the database by their ID. It returns the user and an error if any.
- UpdateUser updates the user's information in the database. It returns the updated user and an error if any.

2. UpdateUser resolves the updateUser mutation, allowing an authenticated user to update their profile information.
2.1 Step 1: Retrieve the user ID from the context (set by the authentication middleware)
2.2 etch the user from the database to ensure they exist
2.3 Validate and prepare the new values for update
2.3 Update the user in the database

## delete user
1. user_repository
DeleteUser deletes a user from the database by their ID. It returns a boolean indicating whether the user was deleted and an error if any.

2. Mutations
DeleteUser resolves the deleteUser mutation, allowing an authenticated user to delete their account.

## get user
>internal/graph/resolvers/query.go
- Me is the resolver for the me field.

## Create habit
1. habit_repository
- CreateHabit creates a new habit for a user in the database with the provided name and description.

2. internal/grapsh/resolvers/mutation.go
CreateHabit is the resolver for the createHabit field. 
**main.go**
3. Initialize repositories
	habitRepo := repository.NewHabitRepository(db)
    4. Wire habit repository to the app
    		Resolvers: &resolvers.Resolver{
			HabitRepo:    habitRepo,
		},

## Get habits
1. habit_repository
- GetHabitsByUserID retrieves all habits associated with a specific user ID from the database.
2. internal/grapsh/resolvers/query.go
Habits is the resolver for the habits field.

## Get habit by ID
1. habit_repository
- GetHabitByID retrieves a habit from the database by its ID. It returns the habit if found, or nil if not found.
2. internal/grapsh/resolvers/query.go
Habit is the resolver for the habit field.