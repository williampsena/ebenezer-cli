package hyprland

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	cmd "github.com/williampsena/ebenezer-cli/internal/cmd"
	jokes "github.com/williampsena/ebenezer-cli/internal/jokes"
)

var defaultLockMessage = "Powered by hyprlock 🔥"

type HyprlockCmd struct {
	HyprlandCmd
	Dry        bool     `help:"Dry run mode, does not write changes to hyprlock.conf" default:"false"`
	Startup    bool     `help:"Run on startup" default:"false"`
	Message    string   `help:"Message for hyprlock" default:""`
	Jokes      bool     `help:"Use a random joke from icanhazdadjoke or reddit"`
	Provider   []string `help:"Joke providers (icanhazdadjoke, reddit, etc). Can specify multiple." default:"reddit,icanhazdadjoke"`
	ConfigPath string   `help:"Hyprlock default config" default:"$HOME/.config/hypr/hyprlock.conf"`
	Format     string   `help:"Format the message" default:"👉 %s 🤪"`
}

func (w *HyprlockCmd) Run(ctx *cmd.Context) error {
	w.SetupContext(ctx)

	w.ConfigPath = os.ExpandEnv(w.ConfigPath)

	data, err := os.ReadFile(w.ConfigPath)
	if err != nil {
		w.Logger.Error("Error reading hyprlock.conf", "err", err)
		return fmt.Errorf("error while trying to read file hyprlock.conf: %v", err)
	}

	content := string(data)

	re := regexp.MustCompile(`(?ms)(label\s*{[^}]*?\btext\s*=\s*)[^\n]+`)

	message, err := w.getMessage()
	if err != nil {
		w.Logger.Error("Error getting message for hyprlock", "err", err)
		return err
	}

	w.Logger.Debug("Raw message for hyprlock: %s", message)
	message = jokes.ParseJokeHtml(message)

	if w.Format != "" {
		message = jokes.ApplyFormat(message, w.Format)
	}

	if w.Dry {
		w.Logger.Debug("Dry run mode enabled, not writing changes to hyprlock.conf")
		return nil
	}

	updated := re.ReplaceAllString(content, fmt.Sprintf(`${1}%s${3}`, message))

	err = os.WriteFile(w.ConfigPath, []byte(updated), 0644)
	if err != nil {
		w.Logger.Error("Error writing hyprlock.conf", "err", err)
		return err
	}

	return nil
}

func (w *HyprlockCmd) getProvider() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return w.Provider[r.Intn(len(w.Provider))]
}

func (w *HyprlockCmd) getMessage() (string, error) {
	if w.Jokes {
		provider := w.getProvider()
		joke, err := w.fetchJokes(provider)
		if err != nil {
			w.Logger.Error("Error fetching joke from %s:", provider, err)
			return defaultLockMessage, nil
		}
		return joke, nil
	}

	if w.Message == "" {
		return defaultLockMessage, nil
	}

	return w.Message, nil
}

func (w *HyprlockCmd) fetchJokes(provider string) (string, error) {
	var joke string
	var err error
	for i := 0; i < 3; i++ {
		fetcher := jokes.BuildJokeFetcher(w.Logger, provider, !w.Startup)
		joke, err = fetcher.FetchJokes()
		if err == nil {
			w.Logger.Debug("Fetched joke from %s: %s", provider, joke)
			return joke, nil
		}

		w.Logger.Debug("error while trying to fetched joke from %s: %s", provider, joke)

		time.Sleep(500 * time.Millisecond)
	}
	return "", fmt.Errorf("failed to fetch joke after 3 attempts: %v", err)
}
