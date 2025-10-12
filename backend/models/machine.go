package models

import (
	"time"

	"gorm.io/gorm"
)

type Machine struct {
	gorm.Model
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	ConfigJSON    string    `json:"config" gorm:"type:json"`
	LastSimulated time.Time `json:"last_simulated"`
	SimulatedRuns int       `json:"simulated_runs"`
}

func (Machine) Tablename() string {
	return "machines"
}
