package hyprland

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	cmd "github.com/williampsena/ebenezer-cli/internal/cmd"
	core "github.com/williampsena/ebenezer-cli/internal/core"
	"github.com/williampsena/ebenezer-cli/internal/hyprland"
	"github.com/williampsena/ebenezer-cli/internal/shell"
	yaml "gopkg.in/yaml.v3"
)

type CronJobs struct {
	Jobs []CronJob `yaml:"jobs"`
}

type CronJob struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Type        string                 `yaml:"type"`
	Command     string                 `yaml:"command"`
	Args        map[string]interface{} `yaml:"args,omitempty"`
	Interval    time.Duration          `yaml:"interval"`
}

type CronCmd struct {
	HyprlandCmd
	Config string `arg:"" help:"Path to the configuration file" default:"~/.config/hypr/cron.yaml"`
}

type DefinedCron map[string]CronHandlerBuilder

func (w *CronCmd) Run(ctx *cmd.Context) error {
	w.SetupContext(ctx)
	return w.SetupCron(ctx)
}

func (w *CronCmd) SetupCron(ctx *cmd.Context) error {
	done := make(chan bool)

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		w.Logger.Error("Failed to create scheduler", "error", err)
		return err
	}

	crons, err := w.parseCrons()

	if err != nil || crons.Jobs == nil || len(crons.Jobs) == 0 {
		w.Logger.Error("No valid cron jobs provided", "error", err)
		return fmt.Errorf("no valid cron jobs provided")
	}

	w.Logger.Info("Setting up cron jobs", "count", len(crons.Jobs))

	definedCron := w.buildDefinedCrons()

	for _, cron := range crons.Jobs {
		_, err := scheduler.NewJob(
			gocron.DurationJob(
				cron.Interval,
			),
			gocron.NewTask(
				w.jobWrapper(
					w.Logger,
					func() {
						w.Logger.Debug("Running cron job", "name", cron.Name, "type", cron.Command, "interval", cron.Interval)
						handler := w.buildCronHandler(definedCron, cron)

						if err := handler(); err != nil {
							w.Logger.Error("Failed to execute cron job", "name", cron.Name, "error", err)
						} else {
							w.Logger.Info("Successfully executed cron job", "name", cron.Name)
						}
					},
				),
			),
		)
		if err != nil {
			w.Logger.Error("Failed to create cron job", "error", err, "cron", cron)
			return fmt.Errorf("failed to create cron job '%s': %w", cron.Name, err)
		}
	}

	if err != nil {
		w.Logger.Error("Failed to create new job", "error", err)
		return err
	}

	scheduler.Start()

	<-done

	err = scheduler.Shutdown()
	if err != nil {
		w.Logger.Error("Failed to shutdown scheduler", "error", err)
		return err
	}

	return nil
}

func (w *CronCmd) parseCrons() (*CronJobs, error) {
	if w.Config == "" {
		return nil, fmt.Errorf("configuration file path is empty")
	}

	w.Config = core.ResolvePath(w.Config)

	w.Logger.Info("Parsing cron jobs from configuration file", "file", w.Config)

	cronJobs, err := w.parseConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file '%s': %w", w.Config, err)
	}

	if cronJobs == nil || len(cronJobs.Jobs) == 0 {
		w.Logger.Warning("No jobs found in configuration file", "file", w.Config)
		return nil, fmt.Errorf("no jobs found in configuration file '%s'", w.Config)
	}

	return cronJobs, nil
}

func (w *CronCmd) parseConfigFile() (*CronJobs, error) {
	configFile, err := os.ReadFile(w.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", w.Config, err)
	}

	var cronJobs CronJobs

	err = yaml.Unmarshal(configFile, &cronJobs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml file: %w", err)
	}

	return &cronJobs, nil
}

type CronHandlerBuilder func(CronJob) CronHandler
type CronHandler func() error

func (w *CronCmd) NoCronHandler() CronHandler {
	return func() error {
		return fmt.Errorf("no cron handler defined")
	}
}

func (w *CronCmd) buildCronHandler(definedCron DefinedCron, cronJob CronJob) CronHandler {
	switch cronJob.Type {
	case "shell":
		return w.runShell(cronJob)
	case "defined":
		if handler, exists := definedCron[cronJob.Command]; exists {
			return handler(cronJob)
		}
	}

	return w.NoCronHandler()
}

func (w *CronCmd) runShell(cronJob CronJob) CronHandler {
	return func() error {
		parts := strings.Fields(cronJob.Command)

		if len(parts) == 0 {
			return fmt.Errorf("empty command")
		}

		command := parts[0]
		args := parts[1:]

		_, err := w.Shell.Run(shell.RunnerExecutionArgs{
			Command: command,
			Args:    args,
		})

		if err != nil {
			w.Logger.Error("Failed to run shell command", "command", cronJob.Command, "error", err)
			return fmt.Errorf("failed to run shell command '%s': %w", cronJob.Command, err)
		}

		w.Logger.Debug("Successfully executed shell command", "command", cronJob.Command)
		return nil
	}
}

func (w *CronCmd) buildDefinedCrons() DefinedCron {
	return map[string]CronHandlerBuilder{
		"$set_random_wallpaper": func(cronJob CronJob) CronHandler {
			return func() error {
				hyprpaper := hyprland.NewHyprpaper(w.Logger, w.Shell)
				hyprpaper.SetWallpaper(core.ResolvePath(cronJob.Args["path"].(string)))
				return nil
			}
		},
		"$update_lock_screen_phrase": func(cronJob CronJob) CronHandler {
			return w.buildHyprlockHandler(cronJob)
		},
	}
}

func (w *CronCmd) buildHyprlockHandler(cronJob CronJob) CronHandler {
	return func() error {
		hyprlockCmd := HyprlockCmd{
			HyprlandCmd: w.HyprlandCmd,
			ConfigPath:  "~/.config/hypr/hyprlock.conf",
			Jokes:       false,
			Message:     "",
			Format:      "👉 %s 🤪",
			Provider:    []string{"reddit", "icanhazdadjoke"},
		}

		if config, ok := cronJob.Args["config"]; ok {
			hyprlockCmd.ConfigPath = core.ResolvePath(config.(string))
		}

		if jokes, ok := cronJob.Args["jokes"]; ok {
			hyprlockCmd.Jokes = jokes.(bool)
		}

		if msg, ok := cronJob.Args["message"]; ok {
			hyprlockCmd.Message = msg.(string)
		}

		if fmtStr, ok := cronJob.Args["format"]; ok {
			hyprlockCmd.Format = fmtStr.(string)
		}

		if provider, ok := cronJob.Args["provider"]; ok {
			if providers, ok := provider.([]string); ok {
				hyprlockCmd.Provider = providers
			}
		}

		return hyprlockCmd.Run(&cmd.Context{})
	}
}

func (w *CronCmd) jobWrapper(logger core.Logger, fn func()) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Cron job failed with error", "error", r)
			}
		}()
		fn()
	}
}
