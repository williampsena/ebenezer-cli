package hyprland

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	cmd "github.com/williampsena/ebenezer-cli/internal/cmd"
	core "github.com/williampsena/ebenezer-cli/internal/core"
	hyprland "github.com/williampsena/ebenezer-cli/internal/hyprland"
	"github.com/williampsena/ebenezer-cli/internal/shell"
)

var SCRIPT_CHANGE_WALLPAPER = `
hyprctl hyprpaper unload all
hyprctl hyprpaper preload "{{.Image}}"
hyprctl hyprpaper wallpaper "{{.Monitor}},{{.Image}}"
`

type HyprpaperCmd struct {
	HyprlandCmd
	hyprpaper   *hyprland.Hyprpaper
	MonitorName string `arg:"" help:"Monitor name to set the wallpaper on" default:""`
	Path        string `help:"Path to a specific wallpaper file or directory" default:"~/Pictures/Wallpapers/Active"`
	Startup     bool   `help:"Run hyprpaper on startup" default:"false"`
}

const (
	configFile = "~/.config/hypr/hyprpaper.conf"
)

func (h *HyprpaperCmd) Run(ctx *cmd.Context) error {
	h.SetupContext(ctx)
	h.hyprpaper = hyprland.NewHyprpaper(h.Logger, h.Shell)

	if err := h.setMonitorName(); err != nil {
		h.Logger.Error("Error setting monitor name", "error", err)
		return err
	}

	wallpaperPath := core.ResolvePath(h.Path)
	configPath := core.ResolvePath(configFile)

	if h.Startup {
		return h.RunStartup(wallpaperPath, configPath)
	}

	return h.RunSetWallpaper(wallpaperPath, configPath)
}

func (h *HyprpaperCmd) RunStartup(wallpaperPath string, configPath string) error {
	configFile, err := h.buildConfig(configPath, wallpaperPath)

	if err != nil {
		h.Logger.Error("Error building configuration file", "error", err)
		return err
	}

	defer configFile.Close()

	if isDirectory := h.isDirectory(wallpaperPath); isDirectory {
		h.setupRandomWallpaper(wallpaperPath, configFile)
	} else {
		h.setupSingleWallpaper(wallpaperPath, configFile)
	}

	h.Logger.Info("Hyprpaper configuration file created successfully", "path", configPath)

	return nil
}

func (h *HyprpaperCmd) RunSetWallpaper(wallpaperPath string, configPath string) error {
	imageFiles, err := h.hyprpaper.FetchRandomImages(wallpaperPath)
	if err != nil {
		return err
	}
	selectedImage := imageFiles[0]

	tmpl, err := template.New("wallpaperScript").Parse(SCRIPT_CHANGE_WALLPAPER)
	if err != nil {
		h.Logger.Error("Error parsing template", "error", err)
		return err
	}

	var scriptBuilder strings.Builder
	err = tmpl.Execute(&scriptBuilder, map[string]string{
		"Image":   selectedImage,
		"Monitor": h.MonitorName,
	})
	if err != nil {
		h.Logger.Error("Error executing template", "error", err)
		return err
	}
	script := scriptBuilder.String()

	_, err = h.Shell.Run(shell.RunnerExecutionArgs{
		Command: "bash",
		Args:    []string{"-c", script},
	})

	if err != nil {
		h.Logger.Error("Error setting wallpaper", "error", err)
		return err
	}

	return nil
}

func (h *HyprpaperCmd) buildConfig(configPath, wallpaperPath string) (*os.File, error) {
	file, err := os.Create(configPath)
	if err != nil {
		h.Logger.Error("Error creating config file", "error", err)
		return nil, err
	}

	_, err = fmt.Fprintf(file, "# Hyprpaper configuration file\n")
	if err != nil {
		h.Logger.Error("Error writing to config file", "error", err)
		return nil, err
	}

	return file, nil
}

func (h *HyprpaperCmd) setupPreloadImage(file *os.File, imageFile string) error {
	_, err := fmt.Fprintf(file, "preload = %s\n", imageFile)
	if err != nil {
		h.Logger.Error("Error writing preload image", "error", err)
		return err
	}

	return nil
}

func (h *HyprpaperCmd) setupRandomWallpaper(wallpaperPath string, configFile *os.File) error {
	imageFiles, err := h.hyprpaper.FetchRandomImages(wallpaperPath)
	if err != nil {
		return err
	}

	selectedImage := imageFiles[0]

	h.setupPreloadImage(configFile, selectedImage)

	_, err = fmt.Fprintf(configFile, "wallpaper = %s,%s\n", h.MonitorName, selectedImage)
	if err != nil {
		fmt.Printf("Erro ao escrever wallpaper: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func (h *HyprpaperCmd) setupSingleWallpaper(wallpaperPath string, configFile *os.File) error {
	h.setupPreloadImage(configFile, wallpaperPath)

	_, err := fmt.Fprintf(configFile, "wallpaper = %s,%s\n", h.MonitorName, wallpaperPath)
	if err != nil {
		h.Logger.Error("Error writing wallpaper option", "error", err)
		return err
	}

	return nil
}

func (h *HyprpaperCmd) setMonitorName() error {
	if h.MonitorName != "" {
		h.Logger.Debug("Using provided monitor name", "name", h.MonitorName)
		return nil
	}

	monitorName, err := h.hyprpaper.GetMonitorName()

	if err != nil {
		h.Logger.Error("Error getting current monitor", "error", err)
		return fmt.Errorf("failed to get current monitor: %w", err)
	}

	h.Logger.Debug("Current monitor", "name", monitorName)

	h.MonitorName = monitorName
	return nil
}

func (h *HyprpaperCmd) isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		h.Logger.Error("Error checking if path is a directory", "path", path, "error", err)
		return false
	}
	return info.IsDir()
}
