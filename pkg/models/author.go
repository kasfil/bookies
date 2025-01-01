// Package models Application structure model
package models

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kasfil/bookies/pkg/database"
	"github.com/kasfil/bookies/pkg/utilities"
)

// IdentifierURI author URI identity binding
type IdentifierURI struct {
	ID string `uri:"id" binding:"required,number,gte=1"`
}

// AuthorBaseModel author model without an ID, case for creating new record
type AuthorBaseModel struct {
	Name      string  `json:"name" binding:"required,validname,gte=3,lte=65"`
	Email     string  `json:"email" binding:"required,email"`
	BirthDate *string `json:"birth_date" binding:"omitempty,datetime=2006-01-02"`
	Bio       *string `json:"bio"`
}

// AuthorDBModel author model for caching database record
type AuthorDBModel struct {
	ID        int          `json:"id" db:"id"`
	Name      string       `json:"name" db:"name"`
	Email     string       `json:"email" db:"email"`
	BirthDate *pgtype.Date `json:"birth_date" db:"birth_date"`
	Bio       *string      `json:"bio" db:"bio"`
	BookTotal uint         `json:"book_total" db:"book_total"`
}

// FetchAuthorDBModel author models to hold multiple authors database record
type FetchAuthorDBModel struct {
	Page        int             `json:"page"`
	Limit       int             `json:"limit"`
	Next        *int            `json:"next"`
	Prev        *int            `json:"prev"`
	RecordTotal int             `json:"record_total"`
	PageTotal   int             `json:"page_total"`
	Data        []AuthorDBModel `json:"data"`
}

// Insert add new author record to the database
func (m *AuthorDBModel) Insert(author *AuthorBaseModel) error {
	// insert query
	query := `INSERT INTO authors (name, email, birth_date, bio)
	VALUES (@name, @email, @birth_date, @bio)
	RETURNING id, name, email, birth_date, bio`

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

	// Run insert mode using named queries
	err = tx.QueryRow(ctx, query, pgx.NamedArgs{
		"name":       author.Name,
		"email":      author.Email,
		"birth_date": author.BirthDate,
		"bio":        author.Bio,
	}).Scan(&m.ID, &m.Name, &m.Email, &m.BirthDate, &m.Bio)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

// Detail get single author by ID
func (m *AuthorDBModel) Detail() error {
	query := `SELECT 
	a.id AS id,
	a.name AS name,
	a.email as email,
	a.birth_date AS birth_date,
	a.bio AS bio,
	COUNT(b.id) AS book_total
	FROM authors a
	LEFT JOIN books b ON a.id = b.author_id
	WHERE a.id = @id
	GROUP BY a.id`

	// Create background context
	ctx := context.Background()
	// Get database connection pool
	db, err := database.GetConnection(ctx)
	if err != nil {
		return err
	}

	// run query
	row, err := db.Conn.Query(ctx, query, pgx.NamedArgs{
		"id": m.ID,
	})
	if err != nil {
		return err
	}

	*m, err = pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[AuthorDBModel])
	if err != nil {
		return err
	}

	return nil
}

// Update update AuthorDBModel from AuthorBaseModel struct
func (m *AuthorDBModel) Update(data *AuthorBaseModel) error {
	query := `UPDATE authors
	SET name = @name,
		email = @email,
		birth_date = @birth_date,
		bio = @bio
	WHERE id = @id
	RETURNING name, email, birth_date, bio`

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

	err = tx.QueryRow(ctx, query, pgx.NamedArgs{
		"name":       data.Name,
		"email":      data.Email,
		"birth_date": data.BirthDate,
		"bio":        data.Bio,
		"id":         m.ID,
	}).Scan(&m.Name, &m.Email, &m.BirthDate, &m.Bio)
	if err != nil {
		// Rollback transaction on error
		tx.Rollback(ctx)
		return err
	}

	return nil
}

// Delete Detele author record from database
func (m *AuthorDBModel) Delete() error {
	query := `DELETE FROM authors WHERE id = $1`

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

// Fetch get authors database record
func (m *FetchAuthorDBModel) Fetch() error {
	query := `SELECT 
	a.id AS id,
	a.name AS name,
	a.email as email,
	a.birth_date AS birth_date,
	a.bio AS bio,
	COUNT(b.id) AS book_total
	FROM authors a
	LEFT JOIN books b ON a.id = b.author_id
	GROUP BY a.id
	ORDER BY a.id DESC
	LIMIT @limit OFFSET @offset`

	// Create background context
	ctx := context.Background()
	// Get database connection pool
	db, err := database.GetConnection(ctx)
	if err != nil {
		return err
	}

	// we need to get all total record first
	err = db.Conn.QueryRow(ctx, "SELECT COUNT(id) AS total FROM authors").Scan(&m.RecordTotal)
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

	rows, err := db.Conn.Query(ctx, query, pgx.NamedArgs{
		"limit":  m.Limit,
		"offset": m.Limit * (m.Page - 1),
	})
	if err != nil {
		return err
	}

	m.Data, err = pgx.CollectRows(rows, pgx.RowToStructByName[AuthorDBModel])
	if err != nil {
		return err
	}

	return nil
}
