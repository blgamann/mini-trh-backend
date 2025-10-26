package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Integration represents an add-on feature for a node
type Integration struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NodeID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"node_id"`
	Node      *Node          `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	Type      string         `gorm:"type:varchar(50);not null;index" json:"type"` // monitoring, backup, explorer
	Status    string         `gorm:"type:varchar(50);not null;index" json:"status"`
	Config    JSONB          `gorm:"type:jsonb" json:"config,omitempty"`
	Info      JSONB          `gorm:"type:jsonb" json:"info,omitempty"` // URL, credentials 등
	LogPath   string         `gorm:"type:varchar(500)" json:"log_path,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the table name
func (Integration) TableName() string {
	return "integrations"
}
