package shell

import (
	"context"
	"os/exec"
	"syscall"
	"time"

	core "github.com/williampsena/ebenezer-cli/internal/core"
)

type RunnerExecutionArgs struct {
	Command   string   `json:"command"`              // command to execute
	Args      []string `json:"args,omitempty"`       // command arguments
	Env       []string `json:"env,omitempty"`        // environment variables to set for the command
	Dir       string   `json:"dir,omitempty"`        // working directory for the command
	Timeout   int      `json:"timeout,omitempty"`    // in seconds
	Setpgid   bool     `json:"setpgid,omitempty"`    // if true, sets cmd.SysProcAttr.Pgid to a new process group ID
	NilStdout bool     `json:"nil_stdout,omitempty"` // if true, sets cmd.Stdout to nil
	NilStderr bool     `json:"nil_stderr,omitempty"` // if true, sets cmd.Stderr to nil
}

// SetDefaults sets default values for RunnerExecutionArgs fields if they are not set.
func (a *RunnerExecutionArgs) SetDefaults() {
	if a.Timeout == 0 {
		a.Timeout = 10
	}
}

type Runner interface {
	// Run executes a command with the provided arguments and environment variables.
	Run(args RunnerExecutionArgs) (string, error)
	// Executes a command and returns its combined standard output and standard error.
	RunCombinedOutput(args RunnerExecutionArgs) (string, error)
	// Executes a command with the provided arguments and environment variables, returning an error if it fails.
	Start(args RunnerExecutionArgs) (int, error)
}

type runnerImpl struct {
	logger core.Logger
}

func NewRunner(logger core.Logger) Runner {
	return &runnerImpl{logger: logger}
}

func (r *runnerImpl) Run(args RunnerExecutionArgs) (string, error) {
	args.SetDefaults()

	ctx, cancel := r.buildContext(args)
	defer cancel()

	cmd := r.buildCmd(ctx, args)

	output, err := cmd.Output()
	if err != nil {
		r.logger.Error("Command execution failed", "error", err)
		return "", err
	}

	r.logger.Debug("Command output", "output", string(output))

	return string(output), nil
}

func (r *runnerImpl) RunCombinedOutput(args RunnerExecutionArgs) (string, error) {
	args.SetDefaults()

	ctx, cancel := r.buildContext(args)
	defer cancel()

	cmd := r.buildCmd(ctx, args)

	output, err := cmd.CombinedOutput()
	if err != nil {
		r.logger.Error("Command execution failed", "error", err)
		return "", err
	}

	r.logger.Debug("Command combined output", "output", string(output))

	return string(output), nil
}

func (r *runnerImpl) Start(args RunnerExecutionArgs) (int, error) {
	args.SetDefaults()

	ctx, cancel := r.buildContext(args)

	cmd := r.buildCmd(ctx, args)

	if err := cmd.Start(); err != nil {
		r.logger.Error("Failed to start command", "error", err)
		return 0, err
	}

	r.logger.Debug("Command started successfully", "command", args.Command)

	go func() {
		defer cancel()

		if err := cmd.Wait(); err != nil {
			r.logger.Error("Command execution failed", "error", err)
		} else {
			r.logger.Debug("Command executed successfully", "command", args.Command)
		}
	}()

	return cmd.Process.Pid, nil
}

func (r *runnerImpl) buildContext(args RunnerExecutionArgs) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(args.Timeout)*time.Second)
}

func (r *runnerImpl) buildCmd(ctx context.Context, args RunnerExecutionArgs) *exec.Cmd {
	// Implementation of the Run method goes here.
	// This is a placeholder implementation.
	r.logger.Debug("Running command", "command", args.Command, "args", args.Args, "env", args.Env, "dir", args.Dir, "timeout", args.Timeout)

	cmd := exec.CommandContext(ctx, args.Command, args.Args...)
	if args.Dir != "" {
		cmd.Dir = args.Dir
	}

	if len(args.Env) > 0 {
		cmd.Env = append(cmd.Env, args.Env...)
	}

	if args.NilStdout {
		cmd.Stdout = nil
	}

	if args.NilStderr {
		cmd.Stderr = nil
	}

	if args.Setpgid {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
	}

	return cmd
}
