package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// CompanyAllowedDomain represents allowed email domains for user registration
type CompanyAllowedDomain struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Domain      string `json:"domain" db:"domain" form:"domain"`
	IsActive    bool   `json:"is_active" db:"is_active" form:"is_active"`
	Description string `json:"description" db:"description" form:"description"`
}

// String returns the JSON representation of the company allowed domain
func (c CompanyAllowedDomain) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// CompanyAllowedDomains is a collection of company allowed domains
type CompanyAllowedDomains []CompanyAllowedDomain

// String returns the JSON representation of the company allowed domains
func (c CompanyAllowedDomains) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" method
func (c *CompanyAllowedDomain) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Domain, Name: "Domain"},
	), nil
}

// IsDomainAllowed checks if a given domain is in the allowed domains list
func IsDomainAllowed(tx *pop.Connection, domain string) (bool, error) {
	var count int
	err := tx.RawQuery("SELECT COUNT(*) FROM company_allowed_domains WHERE domain = ? AND is_active = true", domain).First(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetActiveDomains returns all active allowed domains
func GetActiveDomains(tx *pop.Connection) (CompanyAllowedDomains, error) {
	var domains CompanyAllowedDomains
	err := tx.Where("is_active = true").Order("domain asc").All(&domains)
	return domains, err
}
