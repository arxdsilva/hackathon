package repository

import "github.com/arxdsilva/hackathon/models"

// RepositoryInterface unifies all repository interfaces with namespaced methods
type RepositoryInterface interface {
	// User operations
	UserCount() (int, error)
	UserFindByEmail(email string) (*models.User, error)
	UserFindByID(id interface{}) (*models.User, error)
	UserFindByIDs(ids []interface{}) (*models.Users, error)
	UserGetRecent(limit int) (*models.Users, error)

	// Hackathon operations
	HackathonCount() (int, error)
	HackathonFindByID(id interface{}) (*models.Hackathon, error)
	HackathonFindByOwnerID(ownerID interface{}) (*models.Hackathons, error)
	HackathonGetRecent(limit int) (*models.Hackathons, error)
	HackathonGetActiveWithSchedule() (*models.Hackathons, error)
	HackathonGetActiveHackathonIDs() ([]int, error)

	// Project operations
	ProjectCount() (int, error)
	ProjectCountActive() (int, error)
	ProjectCountPresenting() (int, error)
	ProjectFindByID(id interface{}) (*models.Project, error)
	ProjectFindByHackathonID(hackathonID interface{}) (*models.Projects, error)
	ProjectFindByUserID(userID interface{}) (*models.Projects, error)
	ProjectFindByUserIDWithHackathon(userID interface{}) (*models.Projects, error)
	ProjectFindPresentingByHackathonID(hackathonID interface{}) (*models.Projects, error)
	ProjectFindPresentingFromActiveHackathons() (*models.Projects, error)
	ProjectGetRecent(limit int) (*models.Projects, error)
	ProjectGetFilesByProjectID(projectID interface{}) (*models.Files, error)
	ProjectGetMembershipsByProjectID(projectID interface{}) (*models.ProjectMemberships, error)
	ProjectCountMembershipsByProjectID(projectID interface{}) (int, error)
	ProjectIsUserMemberOfProject(projectID, userID interface{}) (bool, error)

	// Project Membership operations
	ProjectMembershipFindByProjectIDAndUserID(projectID, userID interface{}) (*models.ProjectMembership, error)
	ProjectMembershipCountByProjectID(projectID interface{}) (int, error)
	ProjectMembershipIsUserMember(projectID, userID interface{}) (bool, error)

	// File operations
	FileFindByID(id interface{}) (*models.File, error)
	FileFindAll() (*models.Files, error)
	FileFindAllHackathons() (*models.Hackathons, error)
	FileFindAllProjects() (*models.Projects, error)

	// Company Allowed Domain operations
	CompanyAllowedDomainIsDomainAllowed(domain string) (bool, error)
	CompanyAllowedDomainFindAllActive() (*models.CompanyAllowedDomains, error)
	CompanyAllowedDomainFindAll() (*models.CompanyAllowedDomains, error)
}

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Count() (int, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id interface{}) (*models.User, error)
	FindByIDs(ids []interface{}) (*models.Users, error)
	GetRecent(limit int) (*models.Users, error)
}

// HackathonRepositoryInterface defines the interface for hackathon repository operations
type HackathonRepositoryInterface interface {
	Count() (int, error)
	FindByID(id interface{}) (*models.Hackathon, error)
	FindByOwnerID(ownerID interface{}) (*models.Hackathons, error)
	GetRecent(limit int) (*models.Hackathons, error)
	GetActiveWithSchedule() (*models.Hackathons, error)
	GetActiveHackathonIDs() ([]int, error)
}

// ProjectRepositoryInterface defines the interface for project repository operations
type ProjectRepositoryInterface interface {
	Count() (int, error)
	CountActive() (int, error)
	CountPresenting() (int, error)
	FindByID(id interface{}) (*models.Project, error)
	FindByHackathonID(hackathonID interface{}) (*models.Projects, error)
	FindByUserID(userID interface{}) (*models.Projects, error)
	FindByUserIDWithHackathon(userID interface{}) (*models.Projects, error)
	FindPresentingByHackathonID(hackathonID interface{}) (*models.Projects, error)
	FindPresentingFromActiveHackathons() (*models.Projects, error)
	GetRecent(limit int) (*models.Projects, error)
	GetFilesByProjectID(projectID interface{}) (*models.Files, error)
	GetMembershipsByProjectID(projectID interface{}) (*models.ProjectMemberships, error)
	CountMembershipsByProjectID(projectID interface{}) (int, error)
	IsUserMemberOfProject(projectID, userID interface{}) (bool, error)
}

// ProjectMembershipRepositoryInterface defines the interface for project membership repository operations
type ProjectMembershipRepositoryInterface interface {
	FindByProjectIDAndUserID(projectID, userID interface{}) (*models.ProjectMembership, error)
	CountByProjectID(projectID interface{}) (int, error)
	IsUserMember(projectID, userID interface{}) (bool, error)
}

// FileRepositoryInterface defines the interface for file repository operations
type FileRepositoryInterface interface {
	FindByID(id interface{}) (*models.File, error)
	FindAll() (*models.Files, error)
	FindAllHackathons() (*models.Hackathons, error)
	FindAllProjects() (*models.Projects, error)
}

// CompanyAllowedDomainRepositoryInterface defines the interface for company allowed domain repository operations
type CompanyAllowedDomainRepositoryInterface interface {
	IsDomainAllowed(domain string) (bool, error)
	FindAllActive() (*models.CompanyAllowedDomains, error)
	FindAll() (*models.CompanyAllowedDomains, error)
}
