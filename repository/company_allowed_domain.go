package repository

import (
	"github.com/arxdsilva/hackathon/models"
	"github.com/gobuffalo/pop/v6"
)

// CompanyAllowedDomainRepository handles company allowed domain-related database operations
type CompanyAllowedDomainRepository struct {
	*BaseRepository
}

// NewCompanyAllowedDomainRepository creates a new company allowed domain repository
func NewCompanyAllowedDomainRepository(conn *pop.Connection) *CompanyAllowedDomainRepository {
	return &CompanyAllowedDomainRepository{
		BaseRepository: NewBaseRepository(conn),
	}
}

// IsDomainAllowed checks if a domain is in the allowed domains list
// If no domains are configured, all domains are allowed
func (r *CompanyAllowedDomainRepository) IsDomainAllowed(domain string) (bool, error) {
	// First check if there are any active domains configured
	var totalActive int
	err := r.conn.RawQuery("SELECT COUNT(*) FROM company_allowed_domains WHERE is_active = true").First(&totalActive)
	if err != nil {
		return false, err
	}

	// If no domains are configured, allow all domains
	if totalActive == 0 {
		return true, nil
	}

	// If domains are configured, check if this specific domain is allowed
	var count int
	err = r.conn.RawQuery("SELECT COUNT(*) FROM company_allowed_domains WHERE domain = ? AND is_active = true", domain).First(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAllActive returns all active allowed domains
func (r *CompanyAllowedDomainRepository) FindAllActive() (*models.CompanyAllowedDomains, error) {
	domains := &models.CompanyAllowedDomains{}
	err := r.conn.Where("is_active = true").Order("domain asc").All(domains)
	return domains, err
}

// FindAll returns all domains (active and inactive)
func (r *CompanyAllowedDomainRepository) FindAll() (*models.CompanyAllowedDomains, error) {
	domains := &models.CompanyAllowedDomains{}
	err := r.conn.Order("domain asc").All(domains)
	return domains, err
}
