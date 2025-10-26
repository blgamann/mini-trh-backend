package entities

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONB is a custom type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface.
//
// Called automatically by database/sql when inserting or updating records.
// Converts the Go map (JSONB) into a JSON-encoded []byte so that PostgreSQL
// can store it in a JSONB column.
//
// Flow: Go struct → JSON (as []byte) → PostgreSQL JSONB column
//
// Example:
//
//	user.Profile = JSONB{"name": "Alice"}
//	INSERT INTO users (profile) VALUES ('{"name":"Alice"}');
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface.
//
// Called automatically by database/sql when reading rows from the database.
// Converts the JSONB column (stored as []byte) into a Go map[string]interface{}.
//
// Flow: PostgreSQL JSONB column → JSON (as []byte) → Go map
//
// Example:
//
//	profile column: '{"name":"Alice"}'
//	→ JSONB{"name": "Alice"} in Go
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result
	return nil
}

// Node represents a blockchain node instance
type Node struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Network      string         `gorm:"type:varchar(50);not null" json:"network"` // mainnet, testnet
	Status       string         `gorm:"type:varchar(50);not null;index" json:"status"`
	Config       JSONB          `gorm:"type:jsonb" json:"config,omitempty"`
	Metadata     JSONB          `gorm:"type:jsonb" json:"metadata,omitempty"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Deployments  []Deployment   `gorm:"foreignKey:NodeID" json:"deployments,omitempty"`
	Integrations []Integration  `gorm:"foreignKey:NodeID" json:"integrations,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the default table name used by GORM.
// By default, GORM would use "nodes" (the plural form of "Node"),
// but this method explicitly tells GORM to use "nodes" table.
func (Node) TableName() string {
	return "nodes"
}
