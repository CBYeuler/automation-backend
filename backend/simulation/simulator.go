package simulation

import (
	"log"
	"math/rand"
	"time"

	"github.com/CBYeuler/automation-backend/backend/repository"
)

// MachineSimulator defines the structure to hold dependencies
type MachineSimulator struct {
	Repo repository.MachineRepository
}

// NewMachineSimulator creates a new instance
func NewMachineSimulator(repo repository.MachineRepository) *MachineSimulator {
	return &MachineSimulator{Repo: repo}
}

// StartGlobalSimulation continuously checks for machines and starts/manages simulation goroutines.
func (s *MachineSimulator) StartGlobalSimulation() {
	log.Println("Starting global machine simulation monitor...")

	// This goroutine runs indefinitely to keep the simulation alive and monitor machines
	go func() {
		ticker := time.NewTicker(5 * time.Second) // Check machines every 5 seconds
		defer ticker.Stop()

		// A channel map to track which machine simulation goroutines are running
		// Key: Machine ID, Value: A channel to signal stopping the goroutine
		runningSims := make(map[uint]chan struct{})

		for range ticker.C {
			machines, err := s.Repo.FindAll()
			if err != nil {
				log.Printf("Error fetching machines for simulation: %v", err)
				continue
			}

			// Check for new machines or machines that should be running
			for _, machine := range machines {
				// Only simulate machines with status "Idle" or "Running"
				if (machine.Status == "Idle" || machine.Status == "Running") && runningSims[machine.ID] == nil {
					// Start a new simulation goroutine for this machine
					stopCh := make(chan struct{})
					runningSims[machine.ID] = stopCh
					go s.runMachineSimulation(machine.ID, stopCh)
				}

				// Handle status changes (e.g., if a dashboard command set it to 'Offline')
				if machine.Status == "Offline" && runningSims[machine.ID] != nil {
					// Signal the running goroutine to stop
					close(runningSims[machine.ID])
					delete(runningSims, machine.ID)
					log.Printf("Machine %d (%s) simulation stopped.", machine.ID, machine.Name)
				}
			}
		}
	}()
}

// runMachineSimulation is a long-lived goroutine for a single machine's simulation cycle.
func (s *MachineSimulator) runMachineSimulation(machineID uint, stopCh <-chan struct{}) {
	log.Printf("Machine %d simulation started.", machineID)

	// Update status to Running initially
	s.updateMachineStatus(machineID, "Running")

	// Simulate work cycles
	for {
		select {
		case <-stopCh:
			// Received stop signal
			s.updateMachineStatus(machineID, "Idle") // Set to Idle/Offline upon stopping
			return

		case <-time.After(time.Duration(rand.Intn(4)+1) * time.Second): // Simulate work taking 1-5 seconds
			// Simulation Step
			// machine is a *models.Machine (pointer) because s.Repo.FindByID returns a pointer
			machine, err := s.Repo.FindByID(machineID)
			if err != nil {
				log.Printf("Sim Error: Machine %d not found, stopping simulation.", machineID)
				return // Stop if machine is deleted
			}

			// Core simulation logic: increment runs and update timestamp
			machine.SimulatedRuns++
			machine.LastSimulated = time.Now()

			// Introduce a small chance of error (2%)
			if rand.Intn(100) < 2 {
				s.updateMachineStatus(machineID, "Error")
				log.Printf("Machine %d (%s) has ERROR state!", machineID, machine.Name)
				// Don't return, let the next loop check the status again (e.g., for recovery command)
			} else {
				// Return to Running/Idle if it was in error, or keep Running
				if machine.Status == "Error" {
					s.updateMachineStatus(machineID, "Running")
				}
			}

			if err := s.Repo.Update(machine); err != nil {
				log.Printf("Sim Error: Failed to update machine %d: %v", machineID, err)
			}
			log.Printf("Machine %d (%s) completed run #%d.", machineID, machine.Name, machine.SimulatedRuns)
		}
	}
}

// updateMachineStatus is a helper function to set machine status in DB
func (s *MachineSimulator) updateMachineStatus(machineID uint, status string) {
	// machine is a *models.Machine (pointer)
	machine, err := s.Repo.FindByID(machineID)
	if err != nil {
		log.Printf("Update Status Error: Machine %d not found.", machineID)
		return
	}
	machine.Status = status
	// FIX: Removed '&' since 'machine' is already a pointer
	if err := s.Repo.Update(machine); err != nil {
		log.Printf("Update Status Error: Failed to update status for machine %d: %v", machineID, err)
	}
}
