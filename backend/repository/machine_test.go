package repository_test

import (
	"log"
	"testing"

	"github.com/CBYeuler/automation-backend/backend/models"
	"github.com/CBYeuler/automation-backend/backend/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	// Use an in-memory database connection
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open in-memory DB: %v", err)
	}

	// Migrate the schema (create the table)
	err = db.AutoMigrate(&models.Machine{})
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	return db
}

func TestMachineRepositoryCRUD(t *testing.T) {
	// Setup the test environment
	db := setupTestDB(t)
	repo := repository.NewMachineRepository(db)

	// --- 1. Test Create ---
	t.Run("Create", func(t *testing.T) {
		newMachine := models.Machine{Name: "TestUnit", Status: "Offline", ConfigJSON: `{"temp": 50}`}
		err := repo.Create(&newMachine)

		assert.Nil(t, err, "Create should not return an error")
		assert.Greater(t, newMachine.ID, uint(0), "Machine ID should be set after creation")
	})

	// --- 2. Test FindByID ---
	t.Run("FindByID", func(t *testing.T) {
		machine, err := repo.FindByID(1)

		assert.Nil(t, err, "FindByID should not return an error")
		assert.Equal(t, uint(1), machine.ID, "Retrieved ID should match the created ID")
		assert.Equal(t, "TestUnit", machine.Name, "Retrieved Name should match")
	})

	// --- 3. Test FindAll ---
	t.Run("FindAll", func(t *testing.T) {
		machines, err := repo.FindAll()

		assert.Nil(t, err, "FindAll should not return an error")
		assert.Len(t, machines, 1, "Should find exactly one machine")
	})

	// --- 4. Test Update ---
	t.Run("Update", func(t *testing.T) {
		machine, _ := repo.FindByID(1)
		machine.Status = "Running"
		machine.SimulatedRuns = 10

		err := repo.Update(machine) // Note: The pointer is passed here

		assert.Nil(t, err, "Update should not return an error")

		// Verify update worked by fetching again
		updatedMachine, _ := repo.FindByID(1)
		assert.Equal(t, "Running", updatedMachine.Status, "Status should be updated")
		assert.Equal(t, 10, updatedMachine.SimulatedRuns, "Runs should be updated")
	})

	// --- 5. Test Delete (Soft Delete) ---
	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(1)

		assert.Nil(t, err, "Delete should not return an error")

		// Verify soft delete by trying to FindByID (GORM behavior usually hides soft-deleted records)
		_, err = repo.FindByID(1)
		assert.NotNil(t, err, "Machine should be considered not found after soft delete")
	})
}
