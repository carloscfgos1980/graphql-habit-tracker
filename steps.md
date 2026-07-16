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


goose -dir ./migrations sqlite3 habit.db up
goose -dir ./migrations sqlite3 habit.db status
goose -dir ./migrations sqlite3 habit.db down