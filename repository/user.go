package repository

import (
	"github.com/arxdsilva/hackathon/models"
	"github.com/gobuffalo/pop/v6"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(conn *pop.Connection) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(conn),
	}
}

// Count returns the total number of users
func (r *UserRepository) Count() (int, error) {
	var count int
	err := r.conn.RawQuery("SELECT COUNT(*) FROM users").First(&count)
	return count, err
}

// FindByEmail finds a user by email address
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.conn.Where("email = ?", email).First(user)
	return user, err
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id interface{}) (*models.User, error) {
	user := &models.User{}
	err := r.conn.Find(user, id)
	return user, err
}

// FindByIDs finds multiple users by their IDs
func (r *UserRepository) FindByIDs(ids []interface{}) (*models.Users, error) {
	users := &models.Users{}
	err := r.conn.Where("id IN (?)", ids...).All(users)
	return users, err
}

// GetRecent returns the most recently created users (limited)
func (r *UserRepository) GetRecent(limit int) (*models.Users, error) {
	users := &models.Users{}
	err := r.conn.Order("created_at DESC").Limit(limit).All(users)
	return users, err
}
