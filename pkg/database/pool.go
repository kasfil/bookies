// Package database All in one database connection
package database

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DbPool Database pool structure
type DbPool struct {
	Conn *pgxpool.Pool
}

var (
	dbInstance *DbPool
	dbOnce     sync.Once
)

// Ping only pgxpool Ping wrapper
func (db *DbPool) Ping(ctx context.Context) error {
	return db.Conn.Ping(ctx)
}

// GetConnection Get database connection pool instance, it's also ensuring that
// we only create the connection once
func GetConnection(ctx context.Context) (*DbPool, error) {
	dbOnce.Do(func() {
		connstr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		db, err := pgxpool.New(ctx, connstr)
		if err != nil {
			return
		}

		dbInstance = &DbPool{Conn: db}
	})

	return dbInstance, nil
}
