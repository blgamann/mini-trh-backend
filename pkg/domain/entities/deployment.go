package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deployment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NodeID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"node_id"`
	Node      *Node          `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	Step      string         `gorm:"type:varchar(100);not null" json:"step"` // init, deploy_infra, deploy_chain, finalize
	Status    string         `gorm:"type:varchar(50);not null;index" json:"status"`
	Config    JSONB          `gorm:"type:jsonb" json:"config,omitempty"`
	LogPath   string         `gorm:"type:varchar(500)" json:"log_path,omitempty"`
	ErrorMsg  string         `gorm:"type:text" json:"error_msg,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Deployment) TableName() string {
	return "deployments"
}
