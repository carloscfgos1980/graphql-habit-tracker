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