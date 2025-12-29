package repository

import (
	"github.com/gobuffalo/pop/v6"
)

// Repository defines the base interface for all repositories
type Repository interface {
	GetConnection() *pop.Connection
}

// BaseRepository provides common functionality for all repositories
type BaseRepository struct {
	conn *pop.Connection
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(conn *pop.Connection) *BaseRepository {
	return &BaseRepository{conn: conn}
}

// GetConnection returns the database connection
func (r *BaseRepository) GetConnection() *pop.Connection {
	return r.conn
}
