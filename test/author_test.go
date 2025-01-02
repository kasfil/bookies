package test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/kasfil/bookies/pkg/app"
	"github.com/kasfil/bookies/pkg/database"
	custom_validator "github.com/kasfil/bookies/pkg/validators"
)

var router *gin.Engine

// TestMain setup for file test
func TestMain(m *testing.M) {
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
		validate.RegisterValidation("validname", custom_validator.ValidName)
	}

	router = app.CreateRestApp()

	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	os.Exit(m.Run())
}

// TestFetchAuthors test getting author
func TestFetchAuthors(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/authors", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
