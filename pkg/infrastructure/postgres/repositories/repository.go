package repositories

import "gorm.io/gorm"

// Repositories holds all repository instances
type Repositories struct {
	User        *UserRepository
	Node        *NodeRepository
	Deployment  *DeploymentRepository
	Integration *IntegrationRepository
}

// NewRepositories creates all repository instances
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:        NewUserRepository(db),
		Node:        NewNodeRepository(db),
		Deployment:  NewDeploymentRepository(db),
		Integration: NewIntegrationRepository(db),
	}
}
