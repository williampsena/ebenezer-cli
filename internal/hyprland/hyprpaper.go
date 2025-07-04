package hyprland

import (
	"fmt"
	"html/template"
	"io/fs"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	core "github.com/williampsena/ebenezer-cli/internal/core"
	"github.com/williampsena/ebenezer-cli/internal/shell"
)

var SCRIPT_CHANGE_WALLPAPER = `
hyprctl hyprpaper unload all
hyprctl hyprpaper preload "{{.Image}}"
hyprctl hyprpaper wallpaper "{{.Monitor}},{{.Image}}"
`

type Hyprpaper struct {
	logger core.Logger
	shell  shell.Runner
}

func NewHyprpaper(logger core.Logger, shellRunner shell.Runner) *Hyprpaper {
	return &Hyprpaper{
		logger: logger,
		shell:  shellRunner,
	}
}

func (h *Hyprpaper) SetWallpaper(wallpaperPath string) error {
	monitorName, err := h.GetMonitorName()
	if err != nil {
		h.logger.Error("Error setting monitor name", "error", err)
		return err
	}

	imageFiles, err := h.FetchRandomImages(wallpaperPath)
	if err != nil {
		return err
	}
	selectedImage := imageFiles[0]

	tmpl, err := template.New("wallpaperScript").Parse(SCRIPT_CHANGE_WALLPAPER)
	if err != nil {
		h.logger.Error("Error parsing template", "error", err)
		return err
	}

	var scriptBuilder strings.Builder
	err = tmpl.Execute(&scriptBuilder, map[string]string{
		"Image":   selectedImage,
		"Monitor": monitorName,
	})
	if err != nil {
		h.logger.Error("Error executing template", "error", err)
		return err
	}
	script := scriptBuilder.String()

	_, err = h.shell.Run(shell.RunnerExecutionArgs{
		Command: "bash",
		Args:    []string{"-c", script},
	})

	if err != nil {
		h.logger.Error("Error setting wallpaper", "error", err)
		return err
	}

	return nil
}

func (h *Hyprpaper) FindImageFiles(dir string) ([]string, error) {
	var imageFiles []string
	extensions := []string{".jpg", ".jpeg", ".png"}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, validExt := range extensions {
			if ext == validExt {
				imageFiles = append(imageFiles, path)
				break
			}
		}

		return nil
	})

	return imageFiles, err
}

func (h *Hyprpaper) RandomizeImages(imageFiles []string) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Shuffle(len(imageFiles), func(i, j int) {
		imageFiles[i], imageFiles[j] = imageFiles[j], imageFiles[i]
	})
}

func (h *Hyprpaper) FetchRandomImages(wallpaperPath string) ([]string, error) {
	imageFiles, err := h.FindImageFiles(wallpaperPath)
	if err != nil {
		h.logger.Error("Error finding image files", "error", err)
		return nil, err
	}

	if len(imageFiles) == 0 {
		h.logger.Error("No image files found in the specified directory", "path", wallpaperPath)
		return nil, fmt.Errorf("no image files found in %s", wallpaperPath)
	}

	h.RandomizeImages(imageFiles)

	return imageFiles, nil
}

func (h *Hyprpaper) GetMonitorName() (string, error) {
	output, err := h.shell.Run(shell.RunnerExecutionArgs{
		Command: "hyprctl",
		Args:    []string{"monitors"},
	})
	if err != nil {
		h.logger.Error("Error getting current monitor", "error", err)
		return "", fmt.Errorf("failed to get current monitor: %w", err)
	}
	monitorInfo := string(output)
	if monitorInfo == "" {
		h.logger.Error("No monitors found")
		return "", fmt.Errorf("no monitors found")
	}

	lines := strings.Split(monitorInfo, "\n")
	if len(lines) == 0 {
		h.logger.Error("No monitor information found in output")
		return "", err
	}

	monitorLine := lines[0]
	monitorParts := strings.Fields(monitorLine)

	if len(monitorParts) < 2 {
		h.logger.Error("Invalid monitor information format", "line", monitorLine)
		return "", err
	}

	monitorName := monitorParts[1]
	h.logger.Debug("Current monitor", "name", monitorName)

	return monitorName, nil
}
