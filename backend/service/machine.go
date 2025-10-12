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

func (s *MachineServiceImpl) UpdateMachine(id uint, updatedMachine models.Machine) (models.Machine, error) {
	if updatedMachine.ID != id {
		return models.Machine{}, errors.New("ID mismatch")
	}

	_, err := s.Repo.FindByID(id)
	if err != nil {
		return models.Machine{}, err
	}

	err = s.Repo.Update(&updatedMachine)
	if err != nil {
		return models.Machine{}, err
	}
	return updatedMachine, nil
}

func (s *MachineServiceImpl) DeleteMachine(id uint) error {
	return s.Repo.Delete(id)
}
