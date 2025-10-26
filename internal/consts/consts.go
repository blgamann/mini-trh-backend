package consts

// User Roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Node Status
const (
	NodeStatusPending     = "pending"
	NodeStatusDeploying   = "deploying"
	NodeStatusDeployed    = "deployed"
	NodeStatusFailed      = "failed"
	NodeStatusStopped     = "stopped"
	NodeStatusTerminating = "terminating"
	NodeStatusTerminated  = "terminated"
)

// Deployment Status
const (
	DeploymentStatusPending    = "pending"
	DeploymentStatusInProgress = "in_progress"
	DeploymentStatusCompleted  = "completed"
	DeploymentStatusFailed     = "failed"
	DeploymentStatusStopped    = "stopped"
)

// Integration Types
const (
	IntegrationTypeMonitoring = "monitoring"
	IntegrationTypeBackup     = "backup"
	IntegrationTypeExplorer   = "explorer"
)

// Integration Status
const (
	IntegrationStatusPending    = "pending"
	IntegrationStatusInstalling = "installing"
	IntegrationStatusInstalled  = "installed"
	IntegrationStatusFailed     = "failed"
	IntegrationStatusRemoving   = "removing"
	IntegrationStatusRemoved    = "removed"
)
