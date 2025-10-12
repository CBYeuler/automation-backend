package service

import (
	"errors"

	"github.com/CBYeuler/automation-backend/backend/models"
	"github.com/CBYeuler/automation-backend/backend/repository"
)

type MachineService interface {
	CreateMachine(machine models.Machine) (models.Machine, error)
	GetAllMachines() ([]models.Machine, error)
	GetMachineByID(id uint) (models.Machine, error)
	UpdateMachine(id uint, updatedData models.Machine) (models.Machine, error)
	DeleteMachine(id uint) error
}

type MachineServiceImpl struct {
	Repo repository.MachineRepository
}

func NewMachineService(repo repository.MachineRepository) MachineService {
	return &MachineServiceImpl{Repo: repo}
}

// --- Implementation of the Interface Methods ---
func (s *MachineServiceImpl) CreateMachine(machine models.Machine) (models.Machine, error) {
	if machine.Name == "" {
		return models.Machine{}, errors.New("machine name cannot be empty")
	}
	err := s.Repo.Create(&machine)
	return machine, err
}

func (s *MachineServiceImpl) GetAllMachines() ([]models.Machine, error) {
	return s.Repo.FindAll()
}

func (s *MachineServiceImpl) GetMachineByID(id uint) (models.Machine, error) {
	machine, err := s.Repo.FindByID(id)
	if err != nil {
		return models.Machine{}, err
	}
	return *machine, nil
}

// UpdateMachine handles updates, ensuring the ID is correct and exists.
func (s *MachineServiceImpl) UpdateMachine(id uint, updatedMachine models.Machine) (models.Machine, error) {
	//  Check if the machine exists (important for returning 404, not 500)

	existingMachine, err := s.Repo.FindByID(id)
	if err != nil {
		// Assume gorm.ErrRecordNotFound translates here
		return models.Machine{}, errors.New("machine not found")
	}

	// Enforce the ID from the path (URL parameter)

	updatedMachine.ID = id

	//  Simple copy of fields for demonstration (for full safety, fetch and update field by field)
	existingMachine.Name = updatedMachine.Name
	existingMachine.Status = updatedMachine.Status
	existingMachine.ConfigJSON = updatedMachine.ConfigJSON
	// Note: LastSimulated and SimulatedRuns should be updated by the Simulator, not the API here

	err = s.Repo.Update(existingMachine) // Use the existingMachine pointer after updating its fields
	return *existingMachine, err
}
func (s *MachineServiceImpl) DeleteMachine(id uint) error {
	return s.Repo.Delete(id)
}
