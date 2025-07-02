package process

import (
	"fmt"
	"slices"

	"github.com/williampsena/ebenezer-cli/internal/core"
)

// Stub implementation for testing
type ProcessManagerMockImpl struct {
	logger   core.Logger
	running  []string
	killed   []string
	binaries []string
}

func NewProcessManagerMock(running []string, binaries []string) ProcessManager {
	return &ProcessManagerMockImpl{
		running:  running,
		binaries: binaries,
		killed:   make([]string, 0),
		logger:   core.BuildLogger(false),
	}
}

func (p *ProcessManagerMockImpl) IsProcessRunning(processName string) bool {
	if slices.Contains(p.running, processName) {
		p.logger.Debug("Process is running", "name", processName)
		return true
	}

	return false
}

func (p *ProcessManagerMockImpl) KillProcess(processName string) error {
	if slices.Contains(p.running, processName) {
		p.logger.Debug("Process killed successfully", "name", processName)
		return nil
	}

	p.logger.Warning("Failed to kill process", "name", processName)
	return fmt.Errorf("failed to kill process %s", processName)
}

func (p *ProcessManagerMockImpl) BinaryExists(binary string) (bool, error) {
	if slices.Contains(p.binaries, binary) {
		p.logger.Debug("binary exists", "name", binary)
		return true, nil
	}

	p.logger.Warning("binary does not exist", "name", binary)
	return false, fmt.Errorf("binary %s does not exist", binary)
}

func (p *ProcessManagerMockImpl) SetProcessRunning(processName string) error {
	if !slices.Contains(p.running, processName) {
		p.running = append(p.running, processName)
		p.logger.Debug("Process set to running", "name", processName)
		return nil
	}

	p.logger.Warning("Process already running", "name", processName)
	return fmt.Errorf("process %s is already running", processName)
}
