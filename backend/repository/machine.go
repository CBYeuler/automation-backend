package repository

import (
	"github.com/CBYeuler/automation-backend/backend/models"
	"gorm.io/gorm"
)

// MachineRepository defines the interface for machine data operations
type MachineRepository interface {
	Create(machine *models.Machine) error
	FindAll() ([]models.Machine, error)
	FindByID(id uint) (*models.Machine, error)
	Update(machine *models.Machine) error
	Delete(id uint) error
}

// MachineRepositoryImpl is the concrete implementation of MachineRepository
type MachineRepositoryImpl struct {
	DB *gorm.DB
}

// NewMachineRepository creates a new instance of MachineRepository
func NewMachineRepository(db *gorm.DB) MachineRepository {
	return &MachineRepositoryImpl{DB: db}
}

// --- Implementation of the Interface Methods ---
func (r *MachineRepositoryImpl) Create(machine *models.Machine) error {
	return r.DB.Create(machine).Error
}

func (r *MachineRepositoryImpl) FindAll() ([]models.Machine, error) {
	var machines []models.Machine
	err := r.DB.Find(&machines).Error
	return machines, err
}

func (r *MachineRepositoryImpl) FindByID(id uint) (*models.Machine, error) {
	var machine models.Machine
	err := r.DB.First(&machine, id).Error
	if err != nil {
		return nil, err
	}
	return &machine, nil
}

func (r *MachineRepositoryImpl) Update(machine *models.Machine) error {
	return r.DB.Save(machine).Error
}

func (r *MachineRepositoryImpl) Delete(id uint) error {
	return r.DB.Delete(&models.Machine{}, id).Error
}
