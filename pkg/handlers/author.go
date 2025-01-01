// Package handlers All API handlers
package handlers

import (
	"errors"
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

// AuthorHandler Controllers for author
type AuthorHandler struct{}

// Add insert new author record
func (ac *AuthorHandler) Add(c *gin.Context) {
	var authorBody models.AuthorBaseModel
	if err := c.ShouldBind(&authorBody); err != nil {
		msg := utilities.ParseValidationError(err.(validator.ValidationErrors))
		c.JSON(http.StatusUnprocessableEntity, msg)
		return
	}

	author := new(models.AuthorDBModel)
	err := author.Insert(&authorBody)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				c.JSON(http.StatusConflict, gin.H{"msg": "email already registered"})
			}
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops we made mistake"})
		}
		return
	}

	c.JSON(http.StatusOK, author)
}

// Fetch get list of authors from database
func (ac *AuthorHandler) Fetch(c *gin.Context) {
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

	authors := new(models.FetchAuthorDBModel)
	authors.Page = page
	authors.Limit = limit

	err = authors.Fetch()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		return
	}

	c.JSON(http.StatusOK, authors)
}

// Get get detail author by ID
func (ac *AuthorHandler) Get(c *gin.Context) {
	var authorDetail models.IdentifierURI
	if err := c.ShouldBindUri(&authorDetail); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

	author := new(models.AuthorDBModel)
	author.ID, _ = strconv.Atoi(authorDetail.ID)

	if err := author.Detail(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"msg": "author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		}
		return
	}

	c.JSON(http.StatusOK, author)
}

// Update update author handler by ID
func (ac *AuthorHandler) Update(c *gin.Context) {
	var authorDetail models.IdentifierURI
	if err := c.ShouldBindUri(&authorDetail); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

	var reqBody models.AuthorBaseModel
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

	author := new(models.AuthorDBModel)
	author.ID, _ = strconv.Atoi(authorDetail.ID)

	if err := author.Detail(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"msg": "author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		}
		return
	}

	if err := author.Update(&reqBody); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// Delete remove author by ID handler
func (ac *AuthorHandler) Delete(c *gin.Context) {
	var authorDetail models.IdentifierURI
	if err := c.ShouldBindUri(&authorDetail); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

	author := new(models.AuthorDBModel)
	author.ID, _ = strconv.Atoi(authorDetail.ID)

	if err := author.Detail(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"msg": "author not found"})
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

	c.JSON(http.StatusOK, gin.H{"msg": "Author Removed"})
}

// Books get author books
func (ac *AuthorHandler) Books(c *gin.Context) {
	var authorDetail models.IdentifierURI
	if err := c.ShouldBindUri(&authorDetail); err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			msgs := utilities.ParseValidationError(validatorErr)
			c.JSON(http.StatusUnprocessableEntity, msgs)
			return
		}
	}

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

	authorID, _ := strconv.Atoi(authorDetail.ID)
	books := new(models.FetchBookDBModel)
	books.Page = page
	books.Limit = limit

	if err := books.Fetch(&authorID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "oops, we made a mistake"})
		return
	}

	c.JSON(http.StatusOK, books)
}
