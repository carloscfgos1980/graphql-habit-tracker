package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/database"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/graph/generated"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/graph/resolvers"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/middleware"
	"github.com/carloscfgos1980/graphql-habit-tracker/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file loaded: %v", err)
	}
	// get the database path from environment variable or use default
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/habit.db"
	}
	// Initialize the database
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Printf("Database initialized successfully: %s", dbPath)
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	habitRepo := repository.NewHabitRepository(db)
	habitLogRepo := repository.NewHabitLogRepository(db)
	// Initialize GraphQL server
	graphqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{
			UserRepo:     userRepo,
			HabitRepo:    habitRepo,
			HabitLogRepo: habitLogRepo,
		},
	}))
	// Initialize Gin router
	var router *gin.Engine = gin.Default()
	// Set up routes
	router.SetTrustedProxies(nil)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Habit Streak Tracker is running",
			"status":  "success",
		})
	})
	// GraphQL Playground and GraphQL endpoint
	router.GET("/playground", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/graphql").ServeHTTP(c.Writer, c.Request)
	})
	// Apply authentication middleware to the /graphql endpoint
	router.POST("/graphql", middleware.AuthMiddleware(), func(c *gin.Context) {
		graphqlHandler.ServeHTTP(c.Writer, c.Request)
	})
	// get the port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}
	log.Printf("Starting server on port %s...", port)
	// Start the server and log any errors
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// router.Run() // Start the server on the default port 3000
}
