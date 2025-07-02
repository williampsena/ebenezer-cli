package shell

import (
	"fmt"
	"slices"
	"time"

	core "github.com/williampsena/ebenezer-cli/internal/core"
)

type runnerMockImpl struct {
	logger         core.Logger
	successfulCmds []string
	errorCmds      []string
	timeoutCmds    []string
}

func NewRunnerMock(logger core.Logger, successfulCmds, errorCmds, timeoutCmds []string) Runner {
	return &runnerMockImpl{
		logger:         logger,
		successfulCmds: successfulCmds,
		errorCmds:      errorCmds,
		timeoutCmds:    timeoutCmds,
	}
}
func (r *runnerMockImpl) Run(args RunnerExecutionArgs) (string, error) {
	return r.run(args)
}

func (r *runnerMockImpl) RunCombinedOutput(args RunnerExecutionArgs) (string, error) {
	return r.run(args)
}

func (r *runnerMockImpl) Start(args RunnerExecutionArgs) (int, error) {
	_, err := r.run(args)
	if err != nil {
		r.logger.Error("Command execution failed", "error", err)
		return 0, err
	}
	return 1, nil
}

func (r *runnerMockImpl) run(args RunnerExecutionArgs) (string, error) {
	r.logger.Debug("Mock command execution",
		"command", args.Command,
		"args", args.Args,
		"env", args.Env,
		"dir", args.Dir,
		"timeout", args.Timeout,
	)

	mockOutput := "Mock command output"

	switch {
	case slices.Contains(r.successfulCmds, args.Command):
		r.logger.Debug("Mock command executed successfully", "command", args.Command)
		return mockOutput, nil
	case slices.Contains(r.errorCmds, args.Command):
		r.logger.Error("Mock command execution failed", "command", args.Command)
		return "", fmt.Errorf("mock command execution failed: %s", args.Command)
	case slices.Contains(r.timeoutCmds, args.Command):
		r.logger.Warning("Mock command execution timed out", "command", args.Command)
		if args.Timeout > 0 {
			time.Sleep(time.Duration(args.Timeout) * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
		return "", fmt.Errorf("mock command execution timed out: %s", args.Command)
	default:
		r.logger.Warning("Mock command not recognized", "command", args.Command)
		return "", fmt.Errorf("mock command not recognized: %s", args.Command)
	}
}
