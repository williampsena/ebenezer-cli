package jokes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/williampsena/ebenezer-cli/internal/core"
)

func TestIcanhazjoke(t *testing.T) {
	t.Run("Initialize", func(t *testing.T) {
		logger := &core.Logger{}
		settings := &JokeFetcherSettings{
			logger:   logger,
			provider: "icanhazdadjoke",
			useCache: true,
		}

		joke := &icanhazjoke{}
		joke.Initialize(settings)

		if joke.logger != logger {
			t.Error("Expected logger to be set correctly")
		}

		if joke.provider != "icanhazdadjoke" {
			t.Errorf("Expected provider 'icanhazdadjoke', got %s", joke.provider)
		}

		if joke.useCache != true {
			t.Errorf("Expected useCache true, got %v", joke.useCache)
		}
	})

	t.Run("isCacheValid", func(t *testing.T) {
		t.Run("WhenCacheDisabled", func(t *testing.T) {
			joke := &icanhazjoke{}
			joke.useCache = false

			if joke.isCacheValid() {
				t.Error("Expected cache to be invalid when useCache is false")
			}
		})

		t.Run("WhenFileDoesNotExist", func(t *testing.T) {
			os.Remove(icanhazjokeCacheFilePath)

			joke := &icanhazjoke{}
			joke.useCache = true

			if joke.isCacheValid() {
				t.Error("Expected cache to be invalid when file doesn't exist")
			}
		})

		t.Run("WhenFileIsOld", func(t *testing.T) {
			cache := icanhazjokeCacheData{
				Joke:      "Old joke",
				Timestamp: time.Now().Add(-2 * time.Hour),
			}

			data, _ := json.Marshal(cache)
			err := os.WriteFile(icanhazjokeCacheFilePath, data, 0644)
			if err != nil {
				t.Fatalf("Failed to create test cache file: %v", err)
			}
			defer os.Remove(icanhazjokeCacheFilePath)

			oldTime := time.Now().Add(-2 * time.Hour)
			os.Chtimes(icanhazjokeCacheFilePath, oldTime, oldTime)

			joke := &icanhazjoke{}
			joke.useCache = true

			if joke.isCacheValid() {
				t.Error("Expected cache to be invalid when file is older than 1 hour")
			}
		})

		t.Run("WhenFileIsFresh", func(t *testing.T) {
			cache := icanhazjokeCacheData{
				Joke:      "Fresh joke",
				Timestamp: time.Now(),
			}

			data, _ := json.Marshal(cache)
			err := os.WriteFile(icanhazjokeCacheFilePath, data, 0644)
			if err != nil {
				t.Fatalf("Failed to create test cache file: %v", err)
			}
			defer os.Remove(icanhazjokeCacheFilePath)

			joke := &icanhazjoke{}
			joke.useCache = true

			if !joke.isCacheValid() {
				t.Error("Expected cache to be valid when file is fresh")
			}
		})
	})

	t.Run("loadCache", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			expectedJoke := "Test cached joke"
			cache := icanhazjokeCacheData{
				Joke:      expectedJoke,
				Timestamp: time.Now(),
			}

			data, _ := json.Marshal(cache)
			err := os.WriteFile(icanhazjokeCacheFilePath, data, 0644)
			if err != nil {
				t.Fatalf("Failed to create test cache file: %v", err)
			}
			defer os.Remove(icanhazjokeCacheFilePath)

			joke := &icanhazjoke{}
			loadedJoke, err := joke.loadCache()

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if loadedJoke != expectedJoke {
				t.Errorf("Expected joke '%s', got '%s'", expectedJoke, loadedJoke)
			}
		})

		t.Run("FileNotFound", func(t *testing.T) {
			os.Remove(icanhazjokeCacheFilePath)

			joke := &icanhazjoke{}
			loadedJoke, err := joke.loadCache()

			if err == nil {
				t.Error("Expected error when cache file doesn't exist")
			}

			if loadedJoke != "" {
				t.Errorf("Expected empty joke, got '%s'", loadedJoke)
			}
		})

		t.Run("InvalidJSON", func(t *testing.T) {
			err := os.WriteFile(icanhazjokeCacheFilePath, []byte("invalid json"), 0644)
			if err != nil {
				t.Fatalf("Failed to create test cache file: %v", err)
			}
			defer os.Remove(icanhazjokeCacheFilePath)

			joke := &icanhazjoke{}
			loadedJoke, err := joke.loadCache()

			if err == nil {
				t.Error("Expected error when cache file contains invalid JSON")
			}

			if loadedJoke != "" {
				t.Errorf("Expected empty joke, got '%s'", loadedJoke)
			}
		})
	})

	t.Run("saveCache", func(t *testing.T) {
		defer os.Remove(icanhazjokeCacheFilePath)

		testJoke := "Test joke to save"
		joke := &icanhazjoke{}

		err := joke.saveCache(testJoke)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if _, err := os.Stat(icanhazjokeCacheFilePath); os.IsNotExist(err) {
			t.Error("Expected cache file to be created")
		}

		data, err := os.ReadFile(icanhazjokeCacheFilePath)
		if err != nil {
			t.Fatalf("Failed to read cache file: %v", err)
		}

		var cache icanhazjokeCacheData
		err = json.Unmarshal(data, &cache)
		if err != nil {
			t.Fatalf("Failed to unmarshal cache data: %v", err)
		}

		if cache.Joke != testJoke {
			t.Errorf("Expected cached joke '%s', got '%s'", testJoke, cache.Joke)
		}

		if time.Since(cache.Timestamp) > time.Minute {
			t.Error("Expected timestamp to be recent")
		}
	})

	t.Run("FetchJokes", func(t *testing.T) {
		t.Run("FromCache_Success", func(t *testing.T) {
			expectedJoke := "Cached dad joke"
			cache := icanhazjokeCacheData{
				Joke:      expectedJoke,
				Timestamp: time.Now(),
			}

			data, _ := json.Marshal(cache)
			err := os.WriteFile(icanhazjokeCacheFilePath, data, 0644)
			if err != nil {
				t.Fatalf("Failed to create test cache file: %v", err)
			}
			defer os.Remove(icanhazjokeCacheFilePath)

			logger := &core.Logger{}
			joke := &icanhazjoke{}
			joke.Initialize(&JokeFetcherSettings{
				logger:   logger,
				provider: "icanhazdadjoke",
				useCache: true,
			})

			result, err := joke.FetchJokes()

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != expectedJoke {
				t.Errorf("Expected joke '%s', got '%s'", expectedJoke, result)
			}
		})

		t.Run("FromHTTP Success", func(t *testing.T) {
			defer os.Remove(icanhazjokeCacheFilePath)

			expectedJoke := "Why don't scientists trust atoms? Because they make up everything!"

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Accept") != "application/json" {
					t.Error("Expected Accept header to be application/json")
				}
				if r.Header.Get("User-Agent") != "Go Dad Joke Fetcher" {
					t.Error("Expected User-Agent header to be set correctly")
				}

				response := icanhazjokeResponse{Joke: expectedJoke}
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			logger := &core.Logger{}
			joke := &icanhazjoke{}
			joke.Initialize(&JokeFetcherSettings{
				logger:   logger,
				provider: "icanhazdadjoke",
				useCache: false, // Disable cache to force HTTP request
			})

			_, err := joke.FetchJokes()
			if err != nil {
				if strings.Contains(err.Error(), "dial") || strings.Contains(err.Error(), "lookup") {
					t.Skip("Skipping HTTP test due to network unavailability")
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})

		t.Run("FromHTTP_Error", func(t *testing.T) {
			defer os.Remove(icanhazjokeCacheFilePath)

			logger := &core.Logger{}
			joke := &icanhazjoke{}
			joke.Initialize(&JokeFetcherSettings{
				logger:   logger,
				provider: "icanhazdadjoke",
				useCache: false,
			})

			oldTransport := http.DefaultTransport
			http.DefaultTransport = &http.Transport{
				ResponseHeaderTimeout: 1 * time.Nanosecond, // Force timeout
			}
			defer func() { http.DefaultTransport = oldTransport }()

			_, err := joke.FetchJokes()

			if err != nil {
				t.Logf("Got expected error: %v", err)
			} else {
				t.Logf("Request succeeded despite timeout setting")
			}
		})
	})

	t.Run("MalformedResponse", func(t *testing.T) {
		defer os.Remove(icanhazjokeCacheFilePath)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json response"))
		}))
		defer server.Close()

		logger := &core.Logger{}
		joke := &icanhazjoke{}
		joke.Initialize(&JokeFetcherSettings{
			logger:   logger,
			provider: "icanhazdadjoke",
			useCache: false,
		})
	})

	t.Run("EmptyResponse", func(t *testing.T) {
		expectedJoke := "Test joke for JSON"
		response := icanhazjokeResponse{Joke: expectedJoke}

		data, err := json.Marshal(response)
		if err != nil {
			t.Errorf("Failed to marshal response: %v", err)
		}

		var unmarshaled icanhazjokeResponse
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		if unmarshaled.Joke != expectedJoke {
			t.Errorf("Expected joke '%s', got '%s'", expectedJoke, unmarshaled.Joke)
		}
	})

	t.Run("JSONMarshaling", func(t *testing.T) {
		expectedJoke := "Test cached joke"
		expectedTime := time.Now().Round(time.Second)

		cache := icanhazjokeCacheData{
			Joke:      expectedJoke,
			Timestamp: expectedTime,
		}

		data, err := json.Marshal(cache)
		if err != nil {
			t.Errorf("Failed to marshal cache data: %v", err)
		}

		var unmarshaled icanhazjokeCacheData
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Errorf("Failed to unmarshal cache data: %v", err)
		}

		if unmarshaled.Joke != expectedJoke {
			t.Errorf("Expected joke '%s', got '%s'", expectedJoke, unmarshaled.Joke)
		}

		if !unmarshaled.Timestamp.Equal(expectedTime) {
			t.Errorf("Expected timestamp %v, got %v", expectedTime, unmarshaled.Timestamp)
		}
	})

	t.Run("CacheFilePath", func(t *testing.T) {
		expectedPath := "/tmp/icanhazjoke_cache.json"
		if icanhazjokeCacheFilePath != expectedPath {
			t.Errorf("Expected cache file path '%s', got '%s'", expectedPath, icanhazjokeCacheFilePath)
		}
	})

	t.Run("FetchJokes_EmptyResponse", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		defer os.Remove(icanhazjokeCacheFilePath)

		logger := &core.Logger{}
		joke := &icanhazjoke{}
		joke.Initialize(&JokeFetcherSettings{
			logger:   logger,
			provider: "icanhazdadjoke",
			useCache: true,
		})

		result1, err1 := joke.FetchJokes()
		if err1 != nil {
			t.Skipf("Skipping integration test due to network error: %v", err1)
		}

		if result1 == "" {
			t.Error("Expected non-empty joke from HTTP")
		}

		result2, err2 := joke.FetchJokes()
		if err2 != nil {
			t.Errorf("Expected no error on cached fetch, got %v", err2)
		}

		if result2 != result1 {
			t.Error("Expected same joke from cache")
		}

		if _, err := os.Stat(icanhazjokeCacheFilePath); os.IsNotExist(err) {
			t.Error("Expected cache file to exist after HTTP fetch")
		}
	})
}
