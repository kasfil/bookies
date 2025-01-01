// Package models Application structure model
package models

import (
	"context"
	"math"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kasfil/bookies/pkg/database"
	"github.com/kasfil/bookies/pkg/utilities"
)

// BookBaseModel Book base model, case for creating new record
type BookBaseModel struct {
	Title    string  `json:"title" binding:"required,lte=128,gte=1"`
	Desc     *string `json:"description"`
	PubDate  string  `json:"pub_date" binding:"required,datetime=2006-01-02"`
	AuthorID string  `json:"author_id" binding:"required,number"`
}

// BookDBModel Book database model for structuring database record
type BookDBModel struct {
	ID      int           `json:"id" db:"id"`
	Title   string        `json:"title" db:"title"`
	Desc    *string       `json:"description" db:"description"`
	PubDate *pgtype.Date  `json:"pub_date" db:"publish_date"`
	Author  AuthorDBModel `json:"author" db:"author"`
}

// FetchBookDBModel struct to hold fetch books
type FetchBookDBModel struct {
	Page        int           `json:"page"`
	Limit       int           `json:"limit"`
	Next        *int          `json:"next"`
	Prev        *int          `json:"prev"`
	RecordTotal int           `json:"record_total"`
	PageTotal   int           `json:"page_total"`
	Data        []BookDBModel `json:"data"`
}

// Insert add new book record
func (m *BookDBModel) Insert(data *BookBaseModel) error {
	query := `INSERT INTO books (title, description, publish_date, author_id)
	VALUES (@title, @desc, @pubdate, @author_id)
	RETURNING id, title, description, publish_date, author_id;`

	// Create background context
	ctx := context.Background()
	// Get database connection pool
	db, err := database.GetConnection(ctx)
	if err != nil {
		return err
	}

	// Use transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	// If no error appears, commit transaction
	defer tx.Commit(ctx)

	// Init author db model to attach to the book response
	author := new(AuthorDBModel)

	err = tx.QueryRow(ctx, query, pgx.NamedArgs{
		"title":     data.Title,
		"desc":      data.Desc,
		"pubdate":   data.PubDate,
		"author_id": data.AuthorID,
	}).Scan(&m.ID, &m.Title, &m.Desc, &m.PubDate, &author.ID)
	if err != nil {
		return err
	}

	// Get author detail data
	if err := author.Detail(); err != nil {
		return err
	}

	// set author data
	m.Author = *author

	return nil
}

// Detail get single author by ID
func (m *BookDBModel) Detail() error {
	query := `SELECT
	b.id AS id,
	b.title AS title,
	b.description AS description,
	b.publish_date AS publish_date,
	a.id AS "author.id",
	a.name AS "author.name",
	a.email as "author.email",
	a.birth_date AS "author.birth_date",
	a.bio AS "author.bio",
	(SELECT count(id) FROM books WHERE author_id = a.id) AS "author.book_total"
	FROM books b
	LEFT JOIN authors a ON a.id = b.author_id
	WHERE b.id = $1`

	// Create background context
	ctx := context.Background()
	// Get database connection pool
	db, err := database.GetConnection(ctx)
	if err != nil {
		return err
	}

	// run query
	row, err := db.Conn.Query(ctx, query, m.ID)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanOne(m, row); err != nil {
		return err
	}

	return nil
}

// Update update book record from BookBaseModel struct
func (m *BookDBModel) Update(data *BookBaseModel) error {
	query := `UPDATE books
	SET title = @title,
		description = @desc,
		publish_date = @pub_date,
		author_id = @author_id
	WHERE id = @id
	RETURNING title, description, publish_date, author_id`

	// Create background context
	ctx := context.Background()
	// Get database connection pool
	db, err := database.GetConnection(ctx)
	if err != nil {
		return err
	}

	// Use transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	// Init author db model to attach to the book response
	author := new(AuthorDBModel)

	err = tx.QueryRow(ctx, query, pgx.NamedArgs{
		"title":     data.Title,
		"desc":      data.Desc,
		"pub_date":  data.PubDate,
		"author_id": data.AuthorID,
		"id":        m.ID,
	}).Scan(&m.Title, &m.Desc, &m.PubDate, &author.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	// Get author detail data
	if err := author.Detail(); err != nil {
		tx.Rollback(ctx)
		return err
	}

	// set author data
	m.Author = *author

	return nil
}

// Delete Detele book record from database
func (m *BookDBModel) Delete() error {
	query := `DELETE FROM books WHERE id = $1`

	// Create background context
	ctx := context.Background()
	// Get database connection pool
	db, err := database.GetConnection(ctx)
	if err != nil {
		return err
	}

	// Use transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	result, err := tx.Exec(ctx, query, m.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	} else if result.RowsAffected() > 1 {
		tx.Rollback(ctx)
		return utilities.ErrTooManyAffectedRows
	}

	return nil
}

// Fetch get books database record
func (m *FetchBookDBModel) Fetch(authorID *int) error {
	query := `SELECT
	b.id AS id,
	b.title AS title,
	b.description AS description,
	b.publish_date AS publish_date,
	a.id AS "author.id",
	a.name AS "author.name",
	a.email as "author.email",
	a.birth_date AS "author.birth_date",
	a.bio AS "author.bio",
	(SELECT count(id) FROM books WHERE author_id = a.id) AS "author.book_total"
	FROM books b
	LEFT JOIN authors a ON a.id = b.author_id`

	orderQuery := ` ORDER BY b.id DESC
	LIMIT @limit OFFSET @offset`

	// build query based on authorID existance
	if authorID != nil {
		query = query + " WHERE a.id = @author_id "
	}

	query = query + orderQuery

	// Create background context
	ctx := context.Background()
	// Get database connection pool
	db, err := database.GetConnection(ctx)
	if err != nil {
		return err
	}

	countQuery := "SELECT COUNT(books.id) AS total FROM books"
	if authorID != nil {
		countQuery = countQuery + " JOIN authors ON authors.id = books.author_id where authors.id = @author_id"
	}

	// we need to get all total record first
	err = db.Conn.QueryRow(ctx, countQuery, pgx.NamedArgs{"author_id": authorID}).Scan(&m.RecordTotal)
	if err != nil {
		return err
	}

	// count total page available by divide all record total with limit, if RecordTotal % limit > 1 add one more page
	if m.RecordTotal%m.Limit > 0 {
		m.PageTotal = int(math.Ceil(float64(m.RecordTotal) / float64(m.Limit)))
	} else {
		m.PageTotal = int(m.RecordTotal / m.Limit)
	}

	// Check if user desired page is not higher that available page
	// if so then set to the last page
	if m.Page > m.PageTotal {
		m.Page = m.PageTotal
	}

	// Count m.Prev and m.Next
	if m.Page < 2 {
		m.Prev = nil
	} else {
		prev := m.Page - 1
		m.Prev = &prev
	}

	if m.Page == m.PageTotal {
		m.Next = nil
	} else {
		next := m.Page + 1
		m.Next = &next
	}

	err = pgxscan.Select(ctx, db.Conn, &m.Data, query, pgx.NamedArgs{
		"limit":     m.Limit,
		"offset":    m.Limit * (m.Page - 1),
		"author_id": authorID,
	})
	if err != nil {
		return err
	}

	return nil
}
