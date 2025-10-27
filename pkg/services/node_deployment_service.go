package services

import (
	"errors"
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

var (
	ErrNodeNotFound        = errors.New("node not found")
	ErrDeploymentNotFound  = errors.New("deployment not found")
	ErrIntegrationNotFound = errors.New("integration not found")
	ErrNodeAlreadyExists   = errors.New("node with this name already exists")
)

type NodeDeploymentService struct {
	nodeRepo        *repositories.NodeRepository
	deploymentRepo  *repositories.DeploymentRepository
	integrationRepo *repositories.IntegrationRepository
	taskManager     *taskmanager.TaskManager
}

// NewNodeDeploymentService creates a new NodeDeploymentService instance
func NewNodeDeploymentService(
	nodeRepo *repositories.NodeRepository,
	deploymentRepo *repositories.DeploymentRepository,
	integrationRepo *repositories.IntegrationRepository,
	taskManager *taskmanager.TaskManager,
) *NodeDeploymentService {
	return &NodeDeploymentService{
		nodeRepo:        nodeRepo,
		deploymentRepo:  deploymentRepo,
		integrationRepo: integrationRepo,
		taskManager:     taskManager,
	}
}

func (s *NodeDeploymentService) CreateNode(userID uuid.UUID, name, network string, config entities.JSONB) (*entities.Node, error) {
	exists, err := s.nodeRepo.ExistsByName(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrNodeAlreadyExists
	}

	node := &entities.Node{
		Name:     name,
		Network:  network,
		Status:   consts.NodeStatusPending,
		Config:   config,
		UserID:   userID,
		Metadata: entities.JSONB{},
	}

	if err := s.nodeRepo.Create(node); err != nil {
		logger.Log.Error("Failed to create node", zap.Error(err))
		return nil, err
	}

	deployment := &entities.Deployment{
		NodeID: node.ID,
		Step:   "init",
		Status: consts.DeploymentStatusPending,
		Config: config,
	}
	if err := s.deploymentRepo.Create(deployment); err != nil {
		logger.Log.Error("Failed to create deployment", zap.Error(err))
		return nil, err
	}

	// Add deployment task to queue
	s.taskManager.AddTask(taskmanager.Task{
		Name: fmt.Sprintf("deploy-node-%s", node.ID),
		Fn: func() error {
			return s.deployNode(node.ID, deployment.ID)
		},
	})

	logger.Log.Info("Node creation initiated", zap.String("node_id", node.ID.String()))
	return node, nil
}

// deployNode executes the actual deployment process
func (s *NodeDeploymentService) deployNode(nodeID, deploymentID uuid.UUID) error {
	logger.Log.Info("Starting node deployment", zap.String("node_id", nodeID.String()))

	// Update node status to deploying
	if err := s.nodeRepo.UpdateStatus(nodeID, consts.NodeStatusDeploying); err != nil {
		return err
	}

	// Update deployment status
	if err := s.deploymentRepo.UpdateStatus(deploymentID, consts.DeploymentStatusInProgress); err != nil {
		return err
	}

	// Simulate deployment steps
	steps := []string{"deploy_infrastructure", "deploy_chain", "configure_network", "finalize"}

	for _, step := range steps {
		logger.Log.Info("Executing deployment step",
			zap.String("node_id", nodeID.String()),
			zap.String("step", step),
		)

		// Create deployment record for this step
		stepDeployment := &entities.Deployment{
			NodeID: nodeID,
			Step:   step,
			Status: consts.DeploymentStatusInProgress,
		}
		if err := s.deploymentRepo.Create(stepDeployment); err != nil {
			return s.handleDeploymentFailure(nodeID, deploymentID, err)
		}

		// Simulate step execution
		time.Sleep(2 * time.Second)

		// Update step status
		if err := s.deploymentRepo.UpdateStatus(stepDeployment.ID, consts.DeploymentStatusCompleted); err != nil {
			return s.handleDeploymentFailure(nodeID, deploymentID, err)
		}
	}

	// Mark deployment as completed
	if err := s.deploymentRepo.UpdateStatus(deploymentID, consts.DeploymentStatusCompleted); err != nil {
		return err
	}

	// Update node status to deployed
	if err := s.nodeRepo.UpdateStatus(nodeID, consts.NodeStatusDeployed); err != nil {
		return err
	}

	// Update node metadata with deployment info
	node, _ := s.nodeRepo.FindByID(nodeID)
	if node != nil {
		node.Metadata = entities.JSONB{
			"rpc_url":      "http://node-rpc.example.com",
			"explorer_url": "http://explorer.example.com",
			"deployed_at":  time.Now().Format(time.RFC3339),
		}
		s.nodeRepo.Update(node)
	}

	logger.Log.Info("Node deployment completed successfully", zap.String("node_id", nodeID.String()))
	return nil
}

// handleDeploymentFailure handles deployment failure
func (s *NodeDeploymentService) handleDeploymentFailure(nodeID, deploymentID uuid.UUID, err error) error {
	logger.Log.Error("Deployment failed", zap.String("node_id", nodeID.String()), zap.Error(err))

	// Update deployment status
	s.deploymentRepo.UpdateStatusAndError(deploymentID, consts.DeploymentStatusFailed, err.Error())

	// Update node status
	s.nodeRepo.UpdateStatus(nodeID, consts.NodeStatusFailed)

	return err
}

// GetNode retrieves a node by ID
func (s *NodeDeploymentService) GetNode(nodeID uuid.UUID) (*entities.Node, error) {
	node, err := s.nodeRepo.FindByID(nodeID, "User", "Deployments", "Integrations")
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, ErrNodeNotFound
	}
	return node, nil
}

// ListNodes retrieves all nodes with pagination
func (s *NodeDeploymentService) ListNodes(limit, offset int, filters map[string]interface{}) ([]entities.Node, error) {
	return s.nodeRepo.FindAll(limit, offset, filters)
}

// DeleteNode deletes a node and its related data
func (s *NodeDeploymentService) DeleteNode(nodeID uuid.UUID) error {
	// Check if node exists
	node, err := s.nodeRepo.FindByID(nodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return ErrNodeNotFound
	}

	// Update status to terminating
	if err := s.nodeRepo.UpdateStatus(nodeID, consts.NodeStatusTerminating); err != nil {
		return err
	}

	// Add termination task
	s.taskManager.AddTask(taskmanager.Task{
		Name: fmt.Sprintf("terminate-node-%s", nodeID),
		Fn: func() error {
			return s.terminateNode(nodeID)
		},
	})

	return nil
}

// terminateNode executes the termination process
func (s *NodeDeploymentService) terminateNode(nodeID uuid.UUID) error {
	logger.Log.Info("Terminating node", zap.String("node_id", nodeID.String()))

	// Simulate termination
	time.Sleep(3 * time.Second)

	// Delete related data
	s.integrationRepo.DeleteByNodeID(nodeID)
	s.deploymentRepo.DeleteByNodeID(nodeID)

	// Delete node
	if err := s.nodeRepo.Delete(nodeID); err != nil {
		logger.Log.Error("Failed to delete node", zap.Error(err))
		return err
	}

	logger.Log.Info("Node terminated successfully", zap.String("node_id", nodeID.String()))
	return nil
}

// GetDeployment retrieves a deployment by ID
func (s *NodeDeploymentService) GetDeployment(deploymentID uuid.UUID) (*entities.Deployment, error) {
	deployment, err := s.deploymentRepo.FindByID(deploymentID)
	if err != nil {
		return nil, err
	}
	if deployment == nil {
		return nil, ErrDeploymentNotFound
	}
	return deployment, nil
}

// ListDeployments retrieves all deployments for a node
func (s *NodeDeploymentService) ListDeployments(nodeID uuid.UUID) ([]entities.Deployment, error) {
	return s.deploymentRepo.FindByNodeID(nodeID)
}
