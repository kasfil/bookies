// Package handlers All API controllers
package handlers

import (
	"github.com/gin-gonic/gin"
)

// IncludeHandlers add defined controller to app
func IncludeHandlers(app *gin.Engine) {
	// Author Handler group
	authorC := new(AuthorHandler)
	author := app.Group("/authors")
	author.GET("", authorC.Fetch)
	author.POST("", authorC.Add)
	author.GET("/:id", authorC.Get)
	author.PUT("/:id", authorC.Update)
	author.DELETE("/:id", authorC.Delete)
	author.GET("/:id/books", authorC.Books)

	// Books Handler group
	bookH := new(BookHandler)
	book := app.Group("/books")
	book.GET("", bookH.Fetch)
	book.POST("", bookH.Add)
	book.GET("/:id", bookH.Get)
	book.PUT("/:id", bookH.Update)
	book.DELETE("/:id", bookH.Delete)
}
