// Package dbcontext provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package dbcontext

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// DB represents a DB connection that can be used to run SQL queries.
type DB struct {
	db *dbx.DB
}

type contextKey int

const (
	txKey contextKey = iota
)

// New returns a new DB connection that wraps the given dbx.DB instance.
func New(db *dbx.DB) *DB {
	return &DB{db}
}

// DB returns the dbx.DB wrapped by this object.
func (db *DB) DB() *dbx.DB {
	return db.db
}

// With returns a Builder that can be used to build and execute SQL queries.
// With will return the transaction if it is found in the given context.
// Otherwise it will return a DB connection associated with the context.
func (db *DB) With(ctx context.Context) dbx.Builder {
	if tx, ok := ctx.Value(txKey).(*dbx.Tx); ok {
		return tx
	}
	return db.db.WithContext(ctx)
}