package jokes

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type redditJoke struct{ JokeProvider }

const (
	cacheFilePath = "/tmp/reddit_jokes.json"
)

type cacheData struct {
	Jokes     []string  `json:"jokes"`
	Timestamp time.Time `json:"timestamp"`
}

type redditListing struct {
	Data struct {
		Children []struct {
			Data struct {
				Title    string `json:"title"`
				Selftext string `json:"selftext"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (j *redditJoke) isCacheValid() bool {
	if !j.useCache {
		return false
	}

	info, err := os.Stat(cacheFilePath)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < cacheDuration
}

func (j *redditJoke) loadCache() ([]string, error) {
	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return nil, err
	}

	var cache cacheData
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return nil, err
	}

	return cache.Jokes, nil
}

func (j *redditJoke) saveCache(jokes []string) error {
	cache := cacheData{
		Jokes:     jokes,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFilePath, data, 0644)
}

func (j *redditJoke) FetchJokes() (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if j.isCacheValid() {
		jokes, err := j.loadCache()
		if err == nil && len(jokes) > 0 {
			j.logger.Debug("💾 Using cached jokes")
			return jokes[r.Intn(len(jokes))], nil
		}
	}

	j.logger.Debug("Fetching new jokes from Reddit")

	resp, err := http.Get("https://www.reddit.com/r/ProgrammerDadJokes.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch jokes: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	var listing redditListing
	err = json.Unmarshal(body, &listing)
	if err != nil {
		return "", err
	}

	var jokes []string
	for _, post := range listing.Data.Children {
		title := post.Data.Title
		selftext := post.Data.Selftext
		if selftext != "" {
			jokes = append(jokes, fmt.Sprintf("%s\n%s", title, selftext))
		} else {
			jokes = append(jokes, title)
		}
	}

	if len(jokes) == 0 {
		return "Sorry 🥺, no jokes found", nil
	}

	err = j.saveCache(jokes)

	if err != nil {
		j.logger.Warning("Warning: failed to save cache", "error", err)
		return "", err
	}

	return jokes[r.Intn(len(jokes))], nil
}
