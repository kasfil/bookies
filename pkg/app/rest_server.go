// Package app Provide builder for core application instance
package app

import (
	"github.com/gin-gonic/gin"
	"github.com/kasfil/bookies/pkg/handlers"
)

// CreateRestApp Main rest server builder
func CreateRestApp() *gin.Engine {
	app := gin.Default()

	// include all controllers
	handlers.IncludeHandlers(app)

	return app
}
