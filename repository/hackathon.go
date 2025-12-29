package repository

import (
	"github.com/arxdsilva/hackathon/models"
	"github.com/gobuffalo/pop/v6"
)

// HackathonRepository handles hackathon-related database operations
type HackathonRepository struct {
	*BaseRepository
}

// NewHackathonRepository creates a new hackathon repository
func NewHackathonRepository(conn *pop.Connection) *HackathonRepository {
	return &HackathonRepository{
		BaseRepository: NewBaseRepository(conn),
	}
}

// Count returns the total number of hackathons
func (r *HackathonRepository) Count() (int, error) {
	var count int
	err := r.conn.RawQuery("SELECT COUNT(*) FROM hackathons").First(&count)
	return count, err
}

// FindByID finds a hackathon by ID
func (r *HackathonRepository) FindByID(id interface{}) (*models.Hackathon, error) {
	hackathon := &models.Hackathon{}
	err := r.conn.Find(hackathon, id)
	return hackathon, err
}

// FindByOwnerID finds hackathons owned by a specific user
func (r *HackathonRepository) FindByOwnerID(ownerID interface{}) (*models.Hackathons, error) {
	hackathons := &models.Hackathons{}
	err := r.conn.Where("owner_id = ?", ownerID).All(hackathons)
	return hackathons, err
}

// GetRecent returns the most recently created hackathons (limited)
func (r *HackathonRepository) GetRecent(limit int) (*models.Hackathons, error) {
	hackathons := &models.Hackathons{}
	err := r.conn.Order("created_at DESC").Limit(limit).All(hackathons)
	return hackathons, err
}

// GetActiveWithSchedule returns hackathons that are active/upcoming and have schedules
func (r *HackathonRepository) GetActiveWithSchedule() (*models.Hackathons, error) {
	hackathons := &models.Hackathons{}
	err := r.conn.Where("status IN (?, ?) AND schedule IS NOT NULL AND schedule != ''", "upcoming", "active").Order("start_date asc").All(hackathons)
	return hackathons, err
}

// GetActiveHackathonIDs returns IDs of hackathons with active/upcoming status
func (r *HackathonRepository) GetActiveHackathonIDs() ([]int, error) {
	var ids []int
	err := r.conn.RawQuery("SELECT id FROM hackathons WHERE status IN ('active', 'upcoming')").All(&ids)
	return ids, err
}
