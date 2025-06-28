package widgets

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/williampsena/ebenezer-cli/internal/cmd"
	formatters "github.com/williampsena/ebenezer-cli/internal/cmd/formatters"
)

var distroLogos = map[string]string{
	"ubuntu":        "",
	"fedora":        "",
	"arch":          "",
	"debian":        "",
	"rocky":         "",
	"red-hat":       "",
	"linux-mint":    "",
	"opensuse":      "",
	"manjaro":       "",
	"pop-os":        "",
	"zorin-os":      "",
	"elementary-os": "",
	"solus":         "",
	"void-linux":    "",
	"slackware":     "",
	"artix-linux":   "",
	"endeavouros":   "",
	"other":         "󰌽",
}

type LogoCmd struct {
	WidgetCmd
	Name string `help:"Manually specify the system name." default:""`
	Type string `help:"Output format for the logo." default:"icon" choices:"icon,icon+name,name"`
}

func (w *LogoCmd) Run(ctx *cmd.Context) error {
	w.BuildLogger(ctx.Debug)

	version, err := w.getSystemVersion()
	if err != nil {
		version, err = w.fallbackGetSystemVersion()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching system version: %v\n", err)
			version = "Unknown"
		}
	}

	kernelVersion, err := w.getKernelVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching kernel version: %v\n", err)
		kernelVersion = "Unknown"
	}

	output, err := w.Render(ctx, version, kernelVersion)
	if err != nil {
		return err
	}

	if err := formatters.WriteToStdout(output); err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}

	return nil
}

func (w *LogoCmd) Render(ctx *cmd.Context, version, kernelVersion string) (string, error) {
	distroName := w.Name

	if distroName == "" {
		distroName = version
	}

	distroName = w.parseName(distroName)
	text := w.buildText(distroName)

	data := map[string]interface{}{
		"icon":       w.getLogo(distroName),
		"icon-color": w.IconColor,
		"text":       text,
		"tooltip":    fmt.Sprintf("%s %s", version, kernelVersion),
		"class":      "normal",
		"color":      "#ffffff",
		"no-icon":    w.Type == "name",
	}

	return formatters.FormatWidgetOutput(w.Format, data)
}

func (w *LogoCmd) getSystemVersion() (string, error) {
	file, err := os.Open("/etc/os-release")

	if err != nil {
		return "", err
	}

	defer file.Close()

	var name, version string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NAME=") {
			name = strings.Trim(line[5:], "\"")
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(line[11:], "\"")
		}
	}
	if name == "" {
		name = "Unknown"
	}

	if version != "" {
		return fmt.Sprintf("%s %s", name, version), nil
	}

	return name, nil
}

func (w *LogoCmd) fallbackGetSystemVersion() (string, error) {
	output, err := exec.Command("lsb_release", "-d").Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute lsb_release: %v", err)
	}

	parts := strings.Split(string(output), ":")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected output from lsb_release: %s", output)
	}

	return strings.TrimSpace(parts[1]), nil
}

func (w *LogoCmd) getKernelVersion() (string, error) {
	output, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute uname: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (w *LogoCmd) buildText(distroName string) string {
	switch w.Type {
	case "icon":
		return ""
	case "icon+name":
		return distroName
	case "name":
		return distroName
	default:
		return ""
	}
}

func (w *LogoCmd) getLogo(distroName string) string {
	for key, logo := range distroLogos {
		if strings.HasPrefix(distroName, key) {
			return logo
		}
	}
	return distroLogos["other"]
}

func (w *LogoCmd) parseName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyz0123456789-", r) {
			return r
		}
		return -1
	}, name)
	return name
}
