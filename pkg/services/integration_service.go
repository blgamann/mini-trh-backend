package services

import (
	"fmt"
	"mini-trh-backend/internal/consts"
	"mini-trh-backend/internal/logger"
	"mini-trh-backend/pkg/domain/entities"
	"mini-trh-backend/pkg/infrastructure/postgres/repositories"
	"mini-trh-backend/pkg/services/taskmanager"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IntegrationService struct {
	nodeRepo        *repositories.NodeRepository
	integrationRepo *repositories.IntegrationRepository
	taskManager     *taskmanager.TaskManager
}

func NewIntegrationService(
	nodeRepo *repositories.NodeRepository,
	integrationRepo *repositories.IntegrationRepository,
	taskManager *taskmanager.TaskManager,
) *IntegrationService {
	return &IntegrationService{
		nodeRepo:        nodeRepo,
		integrationRepo: integrationRepo,
		taskManager:     taskManager,
	}
}

// InstallIntegration installs an integration for a node
func (s *IntegrationService) InstallIntegration(nodeID uuid.UUID, integrationType string, config entities.JSONB) (*entities.Integration, error) {
	// Check if node exists
	node, err := s.nodeRepo.FindByID(nodeID)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, ErrNodeNotFound
	}

	// Check if node is deployed
	if node.Status != consts.NodeStatusDeployed {
		return nil, fmt.Errorf("node must be deployed before installing integrations")
	}

	// Check if integration already exists
	exists, err := s.integrationRepo.ExistsByNodeIDAndType(nodeID, integrationType)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("integration %s already installed", integrationType)
	}

	// Create integration record
	integration := &entities.Integration{
		NodeID: nodeID,
		Type:   integrationType,
		Status: consts.IntegrationStatusPending,
		Config: config,
		Info:   entities.JSONB{},
	}

	if err := s.integrationRepo.Create(integration); err != nil {
		logger.Log.Error("Failed to create integration", zap.Error(err))
		return nil, err
	}

	// Add installation task
	s.taskManager.AddTask(taskmanager.Task{
		Name: fmt.Sprintf("install-%s-%s", integrationType, nodeID),
		Fn: func() error {
			return s.installIntegration(integration.ID, nodeID, integrationType)
		},
	})

	logger.Log.Info("Integration installation initiated",
		zap.String("node_id", nodeID.String()),
		zap.String("type", integrationType),
	)
	return integration, nil
}

// installIntegration executes the actual installation
func (s *IntegrationService) installIntegration(integrationID, nodeID uuid.UUID, integrationType string) error {
	logger.Log.Info("Installing integration",
		zap.String("integration_id", integrationID.String()),
		zap.String("type", integrationType),
	)

	// Update status to installing
	if err := s.integrationRepo.UpdateStatus(integrationID, consts.IntegrationStatusInstalling); err != nil {
		return err
	}

	// Simulate installation
	time.Sleep(3 * time.Second)

	// Update integration info with URLs/credentials
	integration, _ := s.integrationRepo.FindByID(integrationID)
	if integration != nil {
		integration.Info = entities.JSONB{
			"url":          fmt.Sprintf("http://%s.node-%s.example.com", integrationType, nodeID),
			"installed_at": time.Now().Format(time.RFC3339),
		}

		// Add type-specific info
		switch integrationType {
		case consts.IntegrationTypeMonitoring:
			integration.Info["dashboard_url"] = fmt.Sprintf("http://grafana.node-%s.example.com", nodeID)
			integration.Info["alert_email"] = ""
		case consts.IntegrationTypeBackup:
			integration.Info["backup_schedule"] = "daily"
			integration.Info["retention_days"] = 30
		case consts.IntegrationTypeExplorer:
			integration.Info["explorer_url"] = fmt.Sprintf("http://explorer.node-%s.example.com", nodeID)
		}

		integration.Status = consts.IntegrationStatusInstalled
		s.integrationRepo.Update(integration)
	}

	logger.Log.Info("Integration installed successfully",
		zap.String("integration_id", integrationID.String()),
		zap.String("type", integrationType),
	)
	return nil
}

// UninstallIntegration removes an integration
func (s *IntegrationService) UninstallIntegration(nodeID uuid.UUID, integrationType string) error {
	// Find integration
	integration, err := s.integrationRepo.FindByNodeIDAndType(nodeID, integrationType)
	if err != nil {
		return err
	}
	if integration == nil {
		return ErrIntegrationNotFound
	}

	// Update status to removing
	if err := s.integrationRepo.UpdateStatus(integration.ID, consts.IntegrationStatusRemoving); err != nil {
		return err
	}

	// Add removal task
	s.taskManager.AddTask(taskmanager.Task{
		Name: fmt.Sprintf("uninstall-%s-%s", integrationType, nodeID),
		Fn: func() error {
			return s.uninstallIntegration(integration.ID, integrationType)
		},
	})

	return nil
}

// uninstallIntegration executes the actual uninstallation
func (s *IntegrationService) uninstallIntegration(integrationID uuid.UUID, integrationType string) error {
	logger.Log.Info("Uninstalling integration",
		zap.String("integration_id", integrationID.String()),
		zap.String("type", integrationType),
	)

	// Simulate uninstallation
	time.Sleep(2 * time.Second)

	// Delete integration
	if err := s.integrationRepo.Delete(integrationID); err != nil {
		logger.Log.Error("Failed to delete integration", zap.Error(err))
		return err
	}

	logger.Log.Info("Integration uninstalled successfully",
		zap.String("integration_id", integrationID.String()),
	)
	return nil
}

// GetIntegration retrieves an integration by ID
func (s *IntegrationService) GetIntegration(integrationID uuid.UUID) (*entities.Integration, error) {
	integration, err := s.integrationRepo.FindByID(integrationID)
	if err != nil {
		return nil, err
	}
	if integration == nil {
		return nil, ErrIntegrationNotFound
	}
	return integration, nil
}

// ListIntegrations retrieves all integrations for a node
func (s *IntegrationService) ListIntegrations(nodeID uuid.UUID) ([]entities.Integration, error) {
	return s.integrationRepo.FindByNodeID(nodeID)
}
