package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CBYeuler/automation-backend/backend/models"
	"github.com/CBYeuler/automation-backend/backend/service"
	"github.com/stretchr/testify/assert"
)

// --- Mock Implementation of the Repository Interface ---

// MockMachineRepository is a mock struct that replaces the real database/repository
type MockMachineRepository struct{}

// Create implements the mock Create method
func (m *MockMachineRepository) Create(machine *models.Machine) error {
	if machine.Name == "ErrorMachine" {
		return errors.New("mock DB error: duplicate name")
	}
	machine.ID = 1 // Simulate successful DB insert
	return nil
}

// FindAll implements the mock FindAll method
// FindAll implements the mock FindAll method
func (m *MockMachineRepository) FindAll() ([]models.Machine, error) {
	return []models.Machine{
		{
			// Corrected: Fully qualify the nested struct: models.Model{...}
			Model:         models.Model{ID: 1},
			Name:          "Assembly Unit 1",
			Status:        "Running",
			ConfigJSON:    "{}",
			LastSimulated: time.Now(),
		},
	}, nil
}

// FindByID implements the mock FindByID method
func (m *MockMachineRepository) FindByID(id uint) (*models.Machine, error) {
	if id == 99 {
		return nil, errors.New("record not found")
	}
	// Corrected: Fully qualify the nested struct: models.Model{...}
	return &models.Machine{
		Model:  models.Model{ID: id},
		Name:   "TestMachine",
		Status: "Idle",
	}, nil
}

// Update implements the mock Update method
func (m *MockMachineRepository) Update(machine *models.Machine) error {
	if machine.ID == 0 {
		return errors.New("mock DB error: update failed (no ID)")
	}
	return nil
}

// Delete implements the mock Delete method
func (m *MockMachineRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("mock DB error: delete failed")
	}
	return nil
}

// --- Test Functions ---

func TestCreateMachineSuccess(t *testing.T) {
	mockRepo := &MockMachineRepository{}
	machineService := service.NewMachineService(mockRepo)

	// Test Case: Valid machine creation
	machine := models.Machine{Name: "NewMachine", Status: "Offline"}
	createdMachine, err := machineService.CreateMachine(machine)

	assert.Nil(t, err, "Error should be nil for successful creation")
	assert.Equal(t, uint(1), createdMachine.ID, "Machine ID should be set by the mock repository")
}

func TestCreateMachineValidationFailure(t *testing.T) {
	mockRepo := &MockMachineRepository{}
	machineService := service.NewMachineService(mockRepo)

	// Test Case: Empty name (Business logic validation)
	machine := models.Machine{Name: "", Status: "Offline"}
	_, err := machineService.CreateMachine(machine)

	assert.NotNil(t, err, "Error should not be nil for validation failure")
	assert.Equal(t, "machine name cannot be empty", err.Error(), "Should return validation error message")
}

func TestGetAllMachines(t *testing.T) {
	mockRepo := &MockMachineRepository{}
	machineService := service.NewMachineService(mockRepo)

	machines, err := machineService.GetAllMachines()

	assert.Nil(t, err, "Error should be nil")
	assert.Len(t, machines, 1, "Should return one mock machine")
	assert.Equal(t, "Assembly Unit 1", machines[0].Name)
}

func TestGetMachineByIDSuccess(t *testing.T) {
	mockRepo := &MockMachineRepository{}
	machineService := service.NewMachineService(mockRepo)

	machine, err := machineService.GetMachineByID(10)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, uint(10), machine.ID, "Should return the correct mock ID")
}

func TestGetMachineByIDNotFound(t *testing.T) {
	mockRepo := &MockMachineRepository{}
	machineService := service.NewMachineService(mockRepo)

	_, err := machineService.GetMachineByID(99)

	assert.NotNil(t, err, "Error should be returned for non-existent ID")
}

func TestUpdateMachineSuccess(t *testing.T) {
	mockRepo := &MockMachineRepository{}
	machineService := service.NewMachineService(mockRepo)

	// Set the ID to a known existing mock ID (10)
	updatedMachine := models.Machine{Model: models.Model{ID: 10}, Name: "UpdatedName", Status: "Running"}

	result, err := machineService.UpdateMachine(10, updatedMachine)

	assert.Nil(t, err, "Error should be nil for successful update")
	assert.Equal(t, uint(10), result.ID, "ID should match the path ID")
	assert.Equal(t, "UpdatedName", result.Name)
}

func TestDeleteMachineSuccess(t *testing.T) {
	mockRepo := &MockMachineRepository{}
	machineService := service.NewMachineService(mockRepo)

	err := machineService.DeleteMachine(1)

	assert.Nil(t, err, "Error should be nil for successful delete")
}
