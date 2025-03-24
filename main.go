package main

import (
	"context"
	"services-api/constants"
	"services-api/db"
	"services-api/handlers"
	"services-api/logger"
	"services-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {

	// initialize logger
	log, err := logger.NewLogger(constants.ServicesLogFileName)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	ctx := context.WithValue(context.Background(), constants.RequestId, uuid.New().String())
	log.Info(ctx, "Starting Services API")

	// Initialize database
	database, err := db.InitDB(ctx, log)
	if err != nil {
		log.Error(ctx, "Failed to connect to database", "error", err)
		panic("failed to connect database: " + err.Error())
	}
	log.Info(ctx, "Database initialized successfully")

	db.SeedData(ctx, database, log)

	// Create repository
	repo := repository.NewServiceRepository(database, log)

	// Set up Gin router with middleware
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		reqCtx := context.WithValue(c.Request.Context(), constants.RequestId, uuid.New().String())
		c.Request = c.Request.WithContext(reqCtx)
		c.Next()
	})

	// Register handlers
	handlers.RegisterServiceHandlers(r, repo, log)

	// Run server
	log.Info(ctx, "Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Error(ctx, "Server failed to start", "error", err)
	}
}
