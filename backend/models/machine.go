package models

import (
	"time"

	"gorm.io/gorm"
)

type Machine struct {
	gorm.Model
	Name          string    `json:"name"`
	status        string    `json:"status"`
	configJSON    string    `json:"config" gorm:"type:json"`
	LastSimulated time.Time `json:"last_simulated"`
	SimulatedRuns int       `json:"simulated_runs"`
}

func (Machine) Tablename() string {
	return "machines"
}
