package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CBYeuler/automation-backend/backend/handler"
	"github.com/CBYeuler/automation-backend/backend/models"
	"github.com/CBYeuler/automation-backend/backend/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// MockMachineRepository is a simple mock for testing the handler/service interaction
type MockMachineRepository struct{}

func (m *MockMachineRepository) Create(machine *models.Machine) error {
	machine.ID = 1
	return nil
}
func (m *MockMachineRepository) FindAll() ([]models.Machine, error) {
	return []models.Machine{
		{Model: models.Model{ID: 1}, Name: "TestMachine", Status: "Idle"},
	}, nil
}
func (m *MockMachineRepository) FindByID(id uint) (*models.Machine, error) {
	if id == 99 {
		return nil, gorm.ErrRecordNotFound // Use gorm error for not found check
	}
	return &models.Machine{Model: models.Model{ID: id}, Name: "TestMachine", Status: "Idle"}, nil
}
func (m *MockMachineRepository) Update(machine *models.Machine) error { return nil }
func (m *MockMachineRepository) Delete(id uint) error                 { return nil }

// setupRouter creates a test router with the handler initialized
func setupRouter() (*gin.Engine, *handler.MachineHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := &MockMachineRepository{}
	machineService := service.NewMachineService(mockRepo)
	machineHandler := handler.NewMachineHandler(machineService)

	// Set up the routes the handler tests will hit
	api := router.Group("/api/v1")
	{
		api.POST("/machines", machineHandler.CreateMachine)
		api.GET("/machines", machineHandler.GetMachines)
		api.GET("/machines/:id", machineHandler.GetMachineByID)
		api.PUT("/machines/:id", machineHandler.UpdateMachine)
		api.DELETE("/machines/:id", machineHandler.DeleteMachine)
	}
	return router, machineHandler
}

func TestCreateMachineHandler(t *testing.T) {
	router, _ := setupRouter()

	// 1. Successful creation
	t.Run("Success", func(t *testing.T) {
		body := `{"name": "New Machine", "status": "Offline", "config_json": "{}"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/machines", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Expected HTTP 201 Created")
		var responseMachine models.Machine
		err := json.Unmarshal(w.Body.Bytes(), &responseMachine)
		assert.Nil(t, err)
		assert.Equal(t, "New Machine", responseMachine.Name)
	})

	// 2. Invalid JSON binding
	t.Run("InvalidInput", func(t *testing.T) {
		body := `{"name": 123, "status": "Offline"}` // Name should be a string
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/machines", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	})
}

func TestGetMachineByIDHandler(t *testing.T) {
	router, _ := setupRouter()

	// 1. Successful retrieval
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/machines/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	})

	// 2. Not Found
	t.Run("NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/machines/99", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	})

	// 3. Invalid ID
	t.Run("InvalidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/machines/abc", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	})
}

// Further tests for GET (All), PUT, and DELETE handlers would follow this pattern.
