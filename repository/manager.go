package repository

import (
	"github.com/arxdsilva/hackathon/models"
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

// User operations
func (rm *RepositoryManager) UserCount() (int, error) {
	return rm.User().Count()
}

func (rm *RepositoryManager) UserFindByEmail(email string) (*models.User, error) {
	return rm.User().FindByEmail(email)
}

func (rm *RepositoryManager) UserFindByID(id interface{}) (*models.User, error) {
	return rm.User().FindByID(id)
}

func (rm *RepositoryManager) UserFindByIDs(ids []interface{}) (*models.Users, error) {
	return rm.User().FindByIDs(ids)
}

func (rm *RepositoryManager) UserGetRecent(limit int) (*models.Users, error) {
	return rm.User().GetRecent(limit)
}

// Hackathon operations
func (rm *RepositoryManager) HackathonCount() (int, error) {
	return rm.Hackathon().Count()
}

func (rm *RepositoryManager) HackathonFindByID(id interface{}) (*models.Hackathon, error) {
	return rm.Hackathon().FindByID(id)
}

func (rm *RepositoryManager) HackathonFindByOwnerID(ownerID interface{}) (*models.Hackathons, error) {
	return rm.Hackathon().FindByOwnerID(ownerID)
}

func (rm *RepositoryManager) HackathonGetRecent(limit int) (*models.Hackathons, error) {
	return rm.Hackathon().GetRecent(limit)
}

func (rm *RepositoryManager) HackathonGetActiveWithSchedule() (*models.Hackathons, error) {
	return rm.Hackathon().GetActiveWithSchedule()
}

func (rm *RepositoryManager) HackathonGetActiveHackathonIDs() ([]int, error) {
	return rm.Hackathon().GetActiveHackathonIDs()
}

// Project operations
func (rm *RepositoryManager) ProjectCount() (int, error) {
	return rm.Project().Count()
}

func (rm *RepositoryManager) ProjectCountActive() (int, error) {
	return rm.Project().CountActive()
}

func (rm *RepositoryManager) ProjectCountPresenting() (int, error) {
	return rm.Project().CountPresenting()
}

func (rm *RepositoryManager) ProjectFindByID(id interface{}) (*models.Project, error) {
	return rm.Project().FindByID(id)
}

func (rm *RepositoryManager) ProjectFindByHackathonID(hackathonID interface{}) (*models.Projects, error) {
	return rm.Project().FindByHackathonID(hackathonID)
}

func (rm *RepositoryManager) ProjectFindByUserID(userID interface{}) (*models.Projects, error) {
	return rm.Project().FindByUserID(userID)
}

func (rm *RepositoryManager) ProjectFindByUserIDWithHackathon(userID interface{}) (*models.Projects, error) {
	return rm.Project().FindByUserIDWithHackathon(userID)
}

func (rm *RepositoryManager) ProjectFindPresentingByHackathonID(hackathonID interface{}) (*models.Projects, error) {
	return rm.Project().FindPresentingByHackathonID(hackathonID)
}

func (rm *RepositoryManager) ProjectFindPresentingFromActiveHackathons() (*models.Projects, error) {
	return rm.Project().FindPresentingFromActiveHackathons()
}

func (rm *RepositoryManager) ProjectGetRecent(limit int) (*models.Projects, error) {
	return rm.Project().GetRecent(limit)
}

func (rm *RepositoryManager) ProjectGetFilesByProjectID(projectID interface{}) (*models.Files, error) {
	return rm.Project().GetFilesByProjectID(projectID)
}

func (rm *RepositoryManager) ProjectGetMembershipsByProjectID(projectID interface{}) (*models.ProjectMemberships, error) {
	return rm.Project().GetMembershipsByProjectID(projectID)
}

func (rm *RepositoryManager) ProjectCountMembershipsByProjectID(projectID interface{}) (int, error) {
	return rm.Project().CountMembershipsByProjectID(projectID)
}

func (rm *RepositoryManager) ProjectIsUserMemberOfProject(projectID, userID interface{}) (bool, error) {
	return rm.Project().IsUserMemberOfProject(projectID, userID)
}

// Project Membership operations
func (rm *RepositoryManager) ProjectMembershipFindByProjectIDAndUserID(projectID, userID interface{}) (*models.ProjectMembership, error) {
	return rm.ProjectMembership().FindByProjectIDAndUserID(projectID, userID)
}

func (rm *RepositoryManager) ProjectMembershipCountByProjectID(projectID interface{}) (int, error) {
	return rm.ProjectMembership().CountByProjectID(projectID)
}

func (rm *RepositoryManager) ProjectMembershipIsUserMember(projectID, userID interface{}) (bool, error) {
	return rm.ProjectMembership().IsUserMember(projectID, userID)
}

// File operations
func (rm *RepositoryManager) FileFindByID(id interface{}) (*models.File, error) {
	return rm.File().FindByID(id)
}

func (rm *RepositoryManager) FileFindAll() (*models.Files, error) {
	return rm.File().FindAll()
}

func (rm *RepositoryManager) FileFindAllHackathons() (*models.Hackathons, error) {
	return rm.File().FindAllHackathons()
}

func (rm *RepositoryManager) FileFindAllProjects() (*models.Projects, error) {
	return rm.File().FindAllProjects()
}

// Company Allowed Domain operations
func (rm *RepositoryManager) CompanyAllowedDomainIsDomainAllowed(domain string) (bool, error) {
	return rm.CompanyAllowedDomain().IsDomainAllowed(domain)
}

func (rm *RepositoryManager) CompanyAllowedDomainFindAllActive() (*models.CompanyAllowedDomains, error) {
	return rm.CompanyAllowedDomain().FindAllActive()
}

func (rm *RepositoryManager) CompanyAllowedDomainFindAll() (*models.CompanyAllowedDomains, error) {
	return rm.CompanyAllowedDomain().FindAll()
}
