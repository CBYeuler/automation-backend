package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/CBYeuler/automation-backend/backend/models"
	"github.com/CBYeuler/automation-backend/backend/service"
	"github.com/gin-gonic/gin"
)

// MachineHandler contains the service interface for dependency injection
type MachineHandler struct {
	Service service.MachineService
}

// NewMachineHandler creates a new handler instance
func NewMachineHandler(s service.MachineService) *MachineHandler {
	return &MachineHandler{Service: s}
}

// CreateMachine handles POST /api/v1/machines
func (h *MachineHandler) CreateMachine(c *gin.Context) {
	var machine models.Machine
	// Bind the incoming JSON request body to the Machine struct
	if err := c.ShouldBindJSON(&machine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdMachine, err := h.Service.CreateMachine(machine)
	if err != nil {
		log.Printf("Error creating machine: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create machine"})
		return
	}

	c.JSON(http.StatusCreated, createdMachine)
}

// GetMachines handles GET /api/v1/machines
func (h *MachineHandler) GetMachines(c *gin.Context) {
	machines, err := h.Service.GetAllMachines()
	if err != nil {
		log.Printf("Error retrieving machines: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve machines"})
		return
	}
	c.JSON(http.StatusOK, machines)
}

// GetMachineByID handles GET /api/v1/machines/:id
func (h *MachineHandler) GetMachineByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
		return
	}

	machine, err := h.Service.GetMachineByID(uint(id))
	if err != nil {
		// Check for gorm.ErrRecordNotFound or similar custom error from service
		c.JSON(http.StatusNotFound, gin.H{"error": "Machine not found"})
		return
	}

	c.JSON(http.StatusOK, machine)
}

// UpdateMachine handles PUT /api/v1/machines/:id
func (h *MachineHandler) UpdateMachine(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
		return
	}

	var machine models.Machine
	if err := c.ShouldBindJSON(&machine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMachine, err := h.Service.UpdateMachine(uint(id), machine)
	if err != nil {
		log.Printf("Error updating machine ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update machine or machine not found"})
		return
	}

	c.JSON(http.StatusOK, updatedMachine)
}

// DeleteMachine handles DELETE /api/v1/machines/:id
func (h *MachineHandler) DeleteMachine(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
		return
	}

	err = h.Service.DeleteMachine(uint(id))
	if err != nil {
		log.Printf("Error deleting machine ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete machine"})
		return
	}

	// Use StatusNoContent for a successful DELETE operation with no body
	c.JSON(http.StatusNoContent, nil)
}
