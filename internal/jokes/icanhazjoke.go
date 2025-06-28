package jokes

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

const (
	icanhazjokeCacheFilePath = "/tmp/icanhazjoke_cache.json"
)

type icanhazjokeCacheData struct {
	Joke      string    `json:"joke"`
	Timestamp time.Time `json:"timestamp"`
}

type icanhazjokeResponse struct {
	Joke string `json:"joke"`
}

type icanhazjoke struct{ JokeProvider }

func (j *icanhazjoke) isCacheValid() bool {
	if !j.useCache {
		return false
	}

	info, err := os.Stat(icanhazjokeCacheFilePath)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < cacheDuration
}

func (j *icanhazjoke) loadCache() (string, error) {
	data, err := os.ReadFile(icanhazjokeCacheFilePath)
	if err != nil {
		return "", err
	}

	var cache icanhazjokeCacheData
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return "", err
	}

	return cache.Joke, nil
}

func (j *icanhazjoke) saveCache(joke string) error {
	cache := icanhazjokeCacheData{
		Joke:      joke,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(icanhazjokeCacheFilePath, data, 0644)
}

func (j *icanhazjoke) FetchJokes() (string, error) {
	if j.isCacheValid() {
		joke, err := j.loadCache()
		if err == nil && joke != "" {
			j.logger.Debug("💾 Using cached jokes")
			return joke, nil
		}
	}

	j.logger.Debug("Fetching new joke from icanhazdadjoke")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Go Dad Joke Fetcher")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var jokeRes icanhazjokeResponse
	err = json.NewDecoder(resp.Body).Decode(&jokeRes)
	if err != nil {
		return "", err
	}

	err = j.saveCache(jokeRes.Joke)
	if err != nil {
		j.logger.Warning("Warning: failed to save cache", "error", err)
		return "", err
	}

	return jokeRes.Joke, nil
}
