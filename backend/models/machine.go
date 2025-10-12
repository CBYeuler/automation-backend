package models

import (
	"time"

	"gorm.io/gorm"
)

// Model is an exported struct that embeds gorm.Model.
// This ensures gorm.Model fields (ID, CreatedAt, etc.) are visible outside the models package.
type Model gorm.Model

// Machine represents a single piece of equipment/machine configuration.
type Machine struct {
	Model             // ⬅️ Use the new exported base model
	Name       string `gorm:"unique;not null" json:"name" binding:"required"`
	Status     string `gorm:"default:'Offline'" json:"status"`
	ConfigJSON string `gorm:"type:jsonb" json:"config_json"`

	// Simulation-specific fields
	LastSimulated time.Time `json:"last_simulated"`
	SimulatedRuns int       `json:"simulated_runs"`
}

// TableName overrides the default table name for better organization
func (Machine) TableName() string {
	return "machines"
}
