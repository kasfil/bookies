// Package handlers All API handlers
package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kasfil/bookies/pkg/models"
	"github.com/kasfil/bookies/pkg/utilities"
)

// BookHandler Controllers for author
type BookHandler struct{}

// Add insert new book record
func (ac *BookHandler) Add(c *gin.Context) {
	var reqBody models.BookBaseModel
	if err := c.ShouldBind(&reqBody); err != nil {
		msg := utilities.ParseValidationError(err.(validator.ValidationErrors))
		c.JSON(http.StatusUnprocessableEntity, msg)
		return
	}

	book := new(models.BookDBModel)
	err := book.Insert(&reqBody)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "unknown author"})
			}
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops we made mistake"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// Fetch get list of book records from database
func (ac *BookHandler) Fetch(c *gin.Context) {
	// Get page value from query params
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "page parameter should be number and greater than 1"})
		return
	}

	// Get limit value from query params
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 5 || limit > 100 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "limit parameter should be number and between 5 and 100"})
		return
	}

	books := new(models.FetchBookDBModel)
	books.Page = page
	books.Limit = limit

	err = books.Fetch(nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// Get get detail book by ID
func (ac *BookHandler) Get(c *gin.Context) {
	var idURI models.IdentifierURI
	if err := c.ShouldBindUri(&idURI); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

	book := new(models.BookDBModel)
	book.ID, _ = strconv.Atoi(idURI.ID)

	if err := book.Detail(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"msg": "book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// Update update book handler by ID
func (ac *BookHandler) Update(c *gin.Context) {
	var idURI models.IdentifierURI
	if err := c.ShouldBindUri(&idURI); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

	var reqBody models.BookBaseModel
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

	book := new(models.BookDBModel)
	book.ID, _ = strconv.Atoi(idURI.ID)

	if err := book.Detail(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"msg": "book not found"})
		} else {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		}
		return
	}

	if err := book.Update(&reqBody); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Delete remove book by ID handler
func (ac *BookHandler) Delete(c *gin.Context) {
	var idURI models.IdentifierURI
	if err := c.ShouldBindUri(&idURI); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

	author := new(models.BookDBModel)
	author.ID, _ = strconv.Atoi(idURI.ID)

	if err := author.Detail(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"msg": "book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		}
		return
	}

	if err := author.Delete(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Book Removed"})
}
