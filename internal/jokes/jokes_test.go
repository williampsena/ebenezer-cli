package jokes

import (
	"strings"
	"testing"
	"time"

	"github.com/williampsena/ebenezer-cli/internal/core"
)

type mockJokeProvider struct {
	shouldReturnError bool
	joke              string
}

func (m *mockJokeProvider) FetchJokes() (string, error) {
	if m.shouldReturnError {
		return "", &mockError{"mock error"}
	}
	return m.joke, nil
}

func (m *mockJokeProvider) Initialize(settings *JokeFetcherSettings) {}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

func TestBuildJokeFetcher(t *testing.T) {
	logger := &core.Logger{}
	provider := "reddit"
	useCache := true

	fetcher := BuildJokeFetcher(logger, provider, useCache)

	if fetcher.settings == nil {
		t.Error("Expected settings to be initialized")
	}

	if fetcher.settings.provider != provider {
		t.Errorf("Expected provider %s, got %s", provider, fetcher.settings.provider)
	}

	if fetcher.settings.useCache != useCache {
		t.Errorf("Expected useCache %v, got %v", useCache, fetcher.settings.useCache)
	}

	if fetcher.settings.logger != logger {
		t.Error("Expected logger to be set correctly")
	}
}

func TestJokeFetcher(t *testing.T) {
	t.Run("FetchJokes", func(t *testing.T) {
		t.Run("Valid Provider", func(t *testing.T) {
			mockProvider := &mockJokeProvider{
				shouldReturnError: false,
				joke:              "Test joke",
			}
			PROVIDERS["mock"] = mockProvider

			defer delete(PROVIDERS, "mock")

			logger := &core.Logger{}
			fetcher := BuildJokeFetcher(logger, "mock", true)

			joke, err := fetcher.FetchJokes()

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if joke != "Test joke" {
				t.Errorf("Expected 'Test joke', got %s", joke)
			}
		})

		t.Run("Invalid Provider", func(t *testing.T) {
			logger := &core.Logger{}
			fetcher := BuildJokeFetcher(logger, "nonexistent", true)

			joke, err := fetcher.FetchJokes()

			if err == nil {
				t.Error("Expected error for nonexistent provider")
			}

			if joke != "" {
				t.Errorf("Expected empty joke, got %s", joke)
			}

			expectedError := "provider nonexistent not found"
			if err.Error() != expectedError {
				t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
			}
		})

		t.Run("MockProviderError", func(t *testing.T) {
			mockProvider := &mockJokeProvider{
				shouldReturnError: true,
				joke:              "",
			}
			PROVIDERS["mock_error"] = mockProvider

			defer delete(PROVIDERS, "mock_error")

			logger := &core.Logger{}
			fetcher := BuildJokeFetcher(logger, "mock_error", true)

			joke, err := fetcher.FetchJokes()

			if err == nil {
				t.Error("Expected error from mock provider")
			}

			if joke != "" {
				t.Errorf("Expected empty joke, got %s", joke)
			}

			if err.Error() != "mock error" {
				t.Errorf("Expected 'mock error', got '%s'", err.Error())
			}
		})
	})
}

func TestParseJokeHtml(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Replace newlines with br tags",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1<br/>Line 2<br/>Line 3",
		},
		{
			name:     "Replace quotes with special character",
			input:    `Say "Hello World"`,
			expected: "Say 󰉾Hello World󰉾",
		},
		{
			name:     "Replace backslashes with spaces",
			input:    "Path\\to\\file",
			expected: "Path to file",
		},
		{
			name:     "Truncate long jokes",
			input:    strings.Repeat("a", 150),
			expected: strings.Repeat("a", 100) + "...",
		},
		{
			name:     "Short joke unchanged",
			input:    "Short joke",
			expected: "Short joke",
		},
		{
			name:     "Exactly 100 characters",
			input:    strings.Repeat("a", 100),
			expected: strings.Repeat("a", 100),
		},
		{
			name:     "Combined transformations",
			input:    "Line 1\nSay \"Hello\" with\\backslash",
			expected: "Line 1<br/>Say 󰉾Hello󰉾 with backslash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseJokeHtml(tt.input)
			if result != tt.expected {
				t.Errorf("ParseJokeHtml() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestApplyFormat(t *testing.T) {
	tests := []struct {
		name     string
		joke     string
		format   string
		expected string
	}{
		{
			name:     "Simple format with %s",
			joke:     "Why did the chicken cross the road?",
			format:   "Joke: %s",
			expected: "Joke: Why did the chicken cross the road?",
		},
		{
			name:     "JSON format",
			joke:     "Test joke",
			format:   "{\"joke\": \"%s\"}",
			expected: "{\"joke\": \"Test joke\"}",
		},
		{
			name:     "No format placeholder",
			joke:     "Test joke",
			format:   "Static text",
			expected: "Static text%!(EXTRA string=Test joke)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyFormat(tt.joke, tt.format)
			if result != tt.expected {
				t.Errorf("ApplyFormat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestJokeProvider_Initialize(t *testing.T) {
	logger := &core.Logger{}
	settings := &JokeFetcherSettings{
		logger:   logger,
		provider: "test",
		useCache: true,
	}

	provider := &JokeProvider{}
	provider.Initialize(settings)

	if provider.JokeFetcherSettings.logger != logger {
		t.Error("Expected logger to be set correctly")
	}

	if provider.JokeFetcherSettings.provider != "test" {
		t.Errorf("Expected provider 'test', got %s", provider.JokeFetcherSettings.provider)
	}

	if provider.JokeFetcherSettings.useCache != true {
		t.Errorf("Expected useCache true, got %v", provider.JokeFetcherSettings.useCache)
	}
}

func TestPROVIDERS_ContainsExpectedProviders(t *testing.T) {
	expectedProviders := []string{"icanhazdadjoke", "reddit"}

	for _, providerName := range expectedProviders {
		if _, exists := PROVIDERS[providerName]; !exists {
			t.Errorf("Expected provider '%s' to be registered", providerName)
		}
	}

	for name, provider := range PROVIDERS {
		if provider == nil {
			t.Errorf("Provider '%s' is nil", name)
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Provider '%s' doesn't implement JokesInterface properly: %v", name, r)
			}
		}()

		logger := &core.Logger{}
		settings := &JokeFetcherSettings{
			logger:   logger,
			provider: name,
			useCache: false,
		}
		provider.Initialize(settings)
	}
}

func TestCacheDuration(t *testing.T) {
	expectedDuration := time.Hour
	if cacheDuration != expectedDuration {
		t.Errorf("Expected cacheDuration to be %v, got %v", expectedDuration, cacheDuration)
	}
}

func TestJokeFetcherSettings_FieldsAccessibility(t *testing.T) {
	logger := &core.Logger{}
	settings := JokeFetcherSettings{
		logger:   logger,
		provider: "test-provider",
		useCache: true,
	}

	if settings.logger != logger {
		t.Error("Logger field not accessible")
	}

	if settings.provider != "test-provider" {
		t.Error("Provider field not accessible")
	}

	if settings.useCache != true {
		t.Error("UseCache field not accessible")
	}
}
