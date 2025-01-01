// Package main provides the entry point for the application.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/kasfil/bookies/pkg/app"
	"github.com/kasfil/bookies/pkg/database"
	"github.com/kasfil/bookies/pkg/validators"
)

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to find .env file")
	}

	ctx := context.Background()

	// Check database connection by create connection pool
	dbconn, err := database.GetConnection(ctx)
	if err != nil {
		log.Fatal("Error create database connection pool", err)
	}

	// try to ping database instance
	err = dbconn.Ping(ctx)
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	// Register custom validator
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("validname", validators.ValidName)
	}

	app := app.CreateRestApp()

	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	app.Run(addr)
}
