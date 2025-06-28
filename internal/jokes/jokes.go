package jokes

import (
	"fmt"
	"strings"
	"time"

	"github.com/williampsena/ebenezer-cli/internal/core"
)

var cacheDuration = time.Hour

type JokesInterface interface {
	FetchJokes() (string, error)
	Initialize(settings *JokeFetcherSettings)
}

var PROVIDERS = map[string]JokesInterface{
	"icanhazdadjoke": &icanhazjoke{},
	"reddit":         &redditJoke{},
}

type JokeProvider struct {
	JokeFetcherSettings
}

func (j *JokeProvider) Initialize(settings *JokeFetcherSettings) {
	j.JokeFetcherSettings = *settings
}

type JokeFetcher struct {
	settings *JokeFetcherSettings
}

type JokeFetcherSettings struct {
	logger   *core.Logger
	provider string
	useCache bool
}

func BuildJokeFetcher(logger *core.Logger, provider string, useCache bool) JokeFetcher {
	return JokeFetcher{
		settings: &JokeFetcherSettings{
			logger:   logger,
			provider: provider,
			useCache: useCache,
		},
	}
}

func (j *JokeFetcher) FetchJokes() (string, error) {
	if provider := PROVIDERS[j.settings.provider]; provider != nil {
		provider.Initialize(j.settings)
		return provider.FetchJokes()
	}

	return "", fmt.Errorf("provider %s not found", j.settings.provider)
}

func ParseJokeHtml(joke string) string {
	joke = strings.ReplaceAll(joke, "\n", "<br/>")
	joke = strings.ReplaceAll(joke, "\"", "󰉾")
	joke = strings.ReplaceAll(joke, "\\", " ")

	if len(joke) > 100 {
		return joke[:100] + "..."
	}

	return joke
}

func ApplyFormat(joke string, format string) string {
	return fmt.Sprintf(format, joke)
}
