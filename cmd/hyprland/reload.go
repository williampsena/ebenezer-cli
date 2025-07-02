package hyprland

import (
	"fmt"
	"os"
	"time"

	cmd "github.com/williampsena/ebenezer-cli/internal/cmd"
	"github.com/williampsena/ebenezer-cli/internal/shell"
)

type ReloadCmd struct {
	HyprlandCmd
	Component string `arg:"" enum:"all,hyprland,waybar" default:"all" help:"Component to reload: 'all', 'hyprland', or 'waybar'"`
	WaitTime  int    `flag:"" short:"w" default:"2" help:"Wait time in seconds between kill and restart"`
}

func (r *ReloadCmd) Run(ctx *cmd.Context) error {
	r.SetupContext(ctx)

	if err := r.validateEnvironment(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	if err := r.checkDependencies(); err != nil {
		r.logger.Warning("Dependency check failed", "error", err)
		return fmt.Errorf("dependency check failed: %w", err)
	}

	r.logger.Info("Starting reload process", "component", r.Component)

	switch r.Component {
	case "all":
		return r.reloadAll()
	case "hyprland":
		return r.reloadHyprland()
	case "waybar":
		return r.reloadWaybar()
	default:
		return fmt.Errorf("invalid component '%s'. Use 'all', 'hyprland', or 'waybar'", r.Component)
	}
}

func (r *ReloadCmd) reloadAll() error {
	r.logger.Info("🔄 Reloading all components")

	if err := r.reloadWaybar(); err != nil {
		r.logger.Warning("Failed to reload waybar", "error", err)
	}

	time.Sleep(1 * time.Second)

	if err := r.reloadHyprland(); err != nil {
		r.logger.Error("❌ Failed to reload hyprland", "error", err)
		return err
	}

	if err := r.performHealthCheck(); err != nil {
		r.logger.Warning("Health check failed after reload", "error", err)
		return err
	}

	r.logger.Info("✅ Successfully reloaded all components")
	return nil
}

func (r *ReloadCmd) reloadHyprland() error {
	r.logger.Info("🔄 Reloading 🔳 Hyprland configuration")

	output, err := r.shell.Run(shell.RunnerExecutionArgs{
		Command: "hyprctl",
		Args:    []string{"reload"},
	})
	if err != nil {
		r.logger.Error("❌ Failed to reload 🔳 Hyprland", "error", err, "output", string(output))
		return fmt.Errorf("failed to reload 🔳 Hyprland: %w", err)
	}

	r.logger.Info("✅ Hyprland 🔳 configuration reloaded successfully")

	if r.Component == "hyprland" {
		if err := r.performHealthCheck(); err != nil {
			r.logger.Warning("🩺 Hyprland 🔳 health check failed", "error", err)
			return err
		}
	}

	return nil
}

func (r *ReloadCmd) reloadWaybar() error {
	r.logger.Info("🔄 Reloading Waybar ➖")

	if !r.isProcessRunning("waybar") {
		r.logger.Info("⏹️ Waybar ➖ is not running, starting it")
		if err := r.startWaybar(); err != nil {
			return err
		}
	} else {
		if err := r.killProcess("waybar"); err != nil {
			r.logger.Warning("Failed to kill Waybar ➖ process", "error", err)
		}

		r.logger.Debug("Waiting for Waybar ➖ to terminate", "seconds", r.WaitTime)
		time.Sleep(time.Duration(r.WaitTime) * time.Second)

		if err := r.startWaybar(); err != nil {
			return err
		}
	}

	if r.Component == "waybar" {
		time.Sleep(1 * time.Second)
		if err := r.performHealthCheck(); err != nil {
			r.logger.Warning("Waybar ➖ health check failed", "error", err)
			return err
		}
	}

	return nil
}

func (r *ReloadCmd) isProcessRunning(processName string) bool {
	return r.processManager.IsProcessRunning(processName)
}

func (r *ReloadCmd) killProcess(processName string) error {
	return r.processManager.KillProcess(processName)
}

func (r *ReloadCmd) startWaybar() error {
	r.logger.Info("🚀 Starting Waybar ➖")

	if exists, _ := r.processManager.BinaryExists("waybar"); !exists {
		return fmt.Errorf("waybar not found in PATH")
	}

	pid, err := r.shell.Start(shell.RunnerExecutionArgs{
		Command:   "waybar",
		Setpgid:   true,
		NilStdout: true,
		NilStderr: true,
	})
	if err != nil {
		r.logger.Error("❌ Failed to start waybar", "error", err)
		return err
	}

	r.logger.Info("✅ Waybar ➖ started successfully", "pid", pid)

	return nil
}

func (r *ReloadCmd) checkDependencies() error {
	dependencies := []string{"hyprctl", "pgrep", "kill"}

	if r.Component == "waybar" || r.Component == "all" {
		dependencies = append(dependencies, "waybar")
	}

	for _, dep := range dependencies {
		if exists, _ := r.processManager.BinaryExists(dep); !exists {
			return fmt.Errorf("dependency '%s' not found in PATH", dep)
		}
	}

	return nil
}

func (r *ReloadCmd) validateEnvironment() error {
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") == "" {
		r.logger.Warning("HYPRLAND_INSTANCE_SIGNATURE not set, may not be in Hyprland session")
	}

	if os.Getenv("WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("not running in Wayland session")
	}

	return nil
}

func (r *ReloadCmd) performHealthCheck() error {
	r.logger.Debug("Performing post-reload health check")

	if r.Component == "hyprland" || r.Component == "all" {
		_, err := r.shell.Run(shell.RunnerExecutionArgs{Command: "hyprctl", Args: []string{"version"}})

		if err != nil {
			r.logger.Warning("❌ Hyprland health check failed", "error", err)
			return fmt.Errorf("hyprland health check failed: %w", err)
		}
		r.logger.Debug("Hyprland health check passed")
	}

	if r.Component == "waybar" || r.Component == "all" {
		if !r.isProcessRunning("waybar") {
			r.logger.Warning("Waybar ➖ is not running after reload")
			return fmt.Errorf("waybar ➖ is not running after reload")
		}
		r.logger.Debug("Waybar ➖ health check passed")
	}

	return nil
}
