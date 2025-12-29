package repository

import (
	"github.com/gobuffalo/pop/v6"
)

// RepositoryManager manages all repositories and provides access to them
type RepositoryManager struct {
	conn *pop.Connection

	userRepo                 *UserRepository
	hackathonRepo            *HackathonRepository
	projectRepo              *ProjectRepository
	projectMembershipRepo    *ProjectMembershipRepository
	fileRepo                 *FileRepository
	companyAllowedDomainRepo *CompanyAllowedDomainRepository
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(conn *pop.Connection) *RepositoryManager {
	return &RepositoryManager{
		conn: conn,
	}
}

// User returns the user repository
func (rm *RepositoryManager) User() *UserRepository {
	if rm.userRepo == nil {
		rm.userRepo = NewUserRepository(rm.conn)
	}
	return rm.userRepo
}

// Hackathon returns the hackathon repository
func (rm *RepositoryManager) Hackathon() *HackathonRepository {
	if rm.hackathonRepo == nil {
		rm.hackathonRepo = NewHackathonRepository(rm.conn)
	}
	return rm.hackathonRepo
}

// Project returns the project repository
func (rm *RepositoryManager) Project() *ProjectRepository {
	if rm.projectRepo == nil {
		rm.projectRepo = NewProjectRepository(rm.conn)
	}
	return rm.projectRepo
}

// ProjectMembership returns the project membership repository
func (rm *RepositoryManager) ProjectMembership() *ProjectMembershipRepository {
	if rm.projectMembershipRepo == nil {
		rm.projectMembershipRepo = NewProjectMembershipRepository(rm.conn)
	}
	return rm.projectMembershipRepo
}

// File returns the file repository
func (rm *RepositoryManager) File() *FileRepository {
	if rm.fileRepo == nil {
		rm.fileRepo = NewFileRepository(rm.conn)
	}
	return rm.fileRepo
}

// CompanyAllowedDomain returns the company allowed domain repository
func (rm *RepositoryManager) CompanyAllowedDomain() *CompanyAllowedDomainRepository {
	if rm.companyAllowedDomainRepo == nil {
		rm.companyAllowedDomainRepo = NewCompanyAllowedDomainRepository(rm.conn)
	}
	return rm.companyAllowedDomainRepo
}
