package hyprland

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	cmd "github.com/williampsena/ebenezer-cli/internal/cmd"
	core "github.com/williampsena/ebenezer-cli/internal/core"
)

type HyprpaperCmd struct {
	HyprlandCmd
	MonitorName string `arg:"" help:"Monitor name to set the wallpaper on" default:"eDP-1"`
	Interval    int    `help:"Interval in seconds for wallpaper slideshow" default:"30"`
	Slideshow   bool   `help:"Enable wallpaper slideshow" default:"false"`
	Path        string `help:"Path to a specific wallpaper file or directory" default:"~/Pictures/Wallpapers/Active"`
}

const (
	configFile = "~/.config/hypr/hyprpaper.conf"
)

func (h *HyprpaperCmd) Run(ctx *cmd.Context) error {
	wallpaperPath := core.ResolvePath(h.Path)
	configPath := core.ResolvePath(configFile)

	configFile, err := h.buildConfig(configPath, wallpaperPath)

	if err != nil {
		h.logger.Error("Error building configuration file", "error", err)
		return err
	}

	defer configFile.Close()

	if h.Slideshow {
		h.setupSlideshow(wallpaperPath, configFile)
	} else {
		h.setupSingleWallpaper(wallpaperPath, configFile)
	}

	h.logger.Info("Hyprpaper configuration file created successfully", "path", configPath)

	return nil
}

func (h *HyprpaperCmd) buildConfig(configPath, wallpaperPath string) (*os.File, error) {
	file, err := os.Create(configPath)
	if err != nil {
		h.logger.Error("Error creating config file", "error", err)
		return nil, err
	}

	_, err = fmt.Fprintf(file, "# Hyprpaper configuration file\n")
	if err != nil {
		h.logger.Error("Error writing to config file", "error", err)
		return nil, err
	}

	_, err = fmt.Fprintf(file, "wallpaper_dir = %s\n", wallpaperPath)
	if err != nil {
		h.logger.Error("Error writing wallpaper directory", "error", err)
		return nil, err
	}

	return file, nil
}

func (h *HyprpaperCmd) findImageFiles(dir string) ([]string, error) {
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

func (h *HyprpaperCmd) setupPreload(file *os.File, imageFiles []string) error {
	for _, img := range imageFiles {
		_, err := fmt.Fprintf(file, "preload = %s\n", img)
		if err != nil {
			h.logger.Error("Error writing preload image", "error", err)
			return err
		}
	}

	return nil
}

func (h *HyprpaperCmd) randomizeImages(imageFiles []string) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Shuffle(len(imageFiles), func(i, j int) {
		imageFiles[i], imageFiles[j] = imageFiles[j], imageFiles[i]
	})
}

func (h *HyprpaperCmd) setupSlideshow(wallpaperPath string, configFile *os.File) error {
	imageFiles, err := h.findImageFiles(wallpaperPath)
	if err != nil {
		h.logger.Error("Error finding image files", "error", err)
		return err
	}

	h.randomizeImages(imageFiles)
	h.setupPreload(configFile, imageFiles)

	if len(imageFiles) == 0 {
		h.logger.Error("No image files found in the specified directory", "path", wallpaperPath)
		return fmt.Errorf("no image files found in %s", wallpaperPath)
	}

	firstImg := imageFiles[0]
	_, err = fmt.Fprintf(configFile, "wallpaper = %s,%s\n", h.MonitorName, firstImg)
	if err != nil {
		fmt.Printf("Erro ao escrever wallpaper: %v\n", err)
		os.Exit(1)
	}

	_, err = fmt.Fprintf(configFile, "slideshow = true\n")
	if err != nil {
		h.logger.Error("Error writing slideshow option", "error", err)
		return err
	}

	_, err = fmt.Fprintf(configFile, "%s", fmt.Sprintf("slideshow_interval = %d\n", h.Interval))
	if err != nil {
		h.logger.Error("Error writing slideshow interval", "error", err)
		return err
	}

	return nil
}

func (h *HyprpaperCmd) setupSingleWallpaper(wallpaperPath string, configFile *os.File) error {
	h.setupPreload(configFile, []string{wallpaperPath})

	_, err := fmt.Fprintf(configFile, "wallpaper = %s,%s\n", h.MonitorName, wallpaperPath)
	if err != nil {
		h.logger.Error("Error writing wallpaper option", "error", err)
		return err
	}

	return nil
}
