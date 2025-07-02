package core

import (
	"fmt"
	"os/exec"
	"strings"
)

// ProcessManager defines an interface for managing system processes.
type ProcessManager interface {
	// isProcessRunning checks if a process with the given name is currently running.
	IsProcessRunning(processName string) bool
	// killProcess attempts to gracefully terminate a process with the given name.
	// If the process does not terminate, it will forcefully kill it.
	KillProcess(processName string) error
}

type processManagerImpl struct {
	logger Logger
}

func NewProcessManager(logger Logger) ProcessManager {
	return &processManagerImpl{
		logger: logger,
	}
}

func (r *processManagerImpl) IsProcessRunning(processName string) bool {
	cmd := exec.Command("pgrep", "-x", processName)
	err := cmd.Run()
	return err == nil
}

func (r *processManagerImpl) KillProcess(processName string) error {
	r.logger.Debug("Killing process", "name", processName)

	cmd := exec.Command("pgrep", "-x", processName)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to find process %s: %w", processName, err)
	}

	pids := strings.Fields(string(output))
	if len(pids) == 0 {
		return fmt.Errorf("no %s processes found", processName)
	}

	for _, pid := range pids {
		r.logger.Debug("Killing PID", "pid", pid)
		killCmd := exec.Command("kill", "-TERM", pid)
		if err := killCmd.Run(); err != nil {
			r.logger.Warning("Failed to kill PID with TERM", "pid", pid, "error", err)
			killCmd = exec.Command("kill", "-KILL", pid)
			if err := killCmd.Run(); err != nil {
				r.logger.Warning("Failed to kill PID with KILL", "pid", pid, "error", err)
			}
		}
	}

	return nil
}
