package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cdipaolo/sentiment"
)

//=======================================================

// Metrics contains the metrics we will output
// for a particular state of a comment thread.
type Metrics struct {
	Comments CommentMetrics `json:"comments"`
	Users    UserMetrics    `json:"users"`
	Actions  ActionMetrics  `json:"actions"`
}

// CommentMetrics contains metrics specific to comments.
type CommentMetrics struct {
	TotalCount    int `json:"total_count"`
	PositiveCount int `json:"positive_count"`
	NegativeCount int `json:"negative_count"`
}

// ActionMetrics contains metrics specific to actions.
type ActionMetrics struct {
	TotalCount int `json:"total_count"`
	LikeCount  int `json:"like_count"`
	FlagCount  int `json:"flag_count"`
}

// UserMetrics contains metrics specific to users.
type UserMetrics struct {
	TotalCount int `json:"total_count"`
}

//======================================================

// Asset defines the structure of an asset in a thread.
type Asset struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Comment defines the structure of a comment
// in a thread.
type Comment struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

// User defines the structure of a user
// in a thread.
type User struct {
	ID   string `json:"id"`
	Name string `json:"displayName"`
}

// Action defines the structure of an action
// in a thread.
type Action struct {
	ID    string `json:"id"`
	Type  string `json:"action_type"`
	Count int    `json:"count"`
}

// Thread defines the structure of a thread.
type Thread struct {
	Assets   []Asset   `json:"assets"`
	Comments []Comment `json:"comments"`
	Users    []User    `json:"users"`
	Actions  []Action  `json:"actions"`
}

//======================================================

const (
	inDir  = "/pfs/stream"
	outDir = "/pfs/out"
)

func main() {

	// Create the sentiment model.
	model, err := sentiment.Restore()
	if err != nil {
		log.Fatal(err)
	}

	if err := filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {

		// Ignore sub-directories.
		if info.IsDir() {
			return nil
		}

		// Otherwise, read the thread data.
		raw, err := ioutil.ReadFile(filepath.Join(inDir, info.Name()))
		if err != nil {
			return err
		}

		// Decode the thread data.
		var thread Thread
		if err = json.Unmarshal(raw, &thread); err != nil {
			return err
		}

		// Get the sentiment for all comments.
		sentimentCounts := make(map[uint8]int)
		for _, comment := range thread.Comments {

			// Get the sentiment for this individual comment.
			analysis := model.SentimentAnalysis(comment.Body, sentiment.English)

			// Increment the respective sentiment counter.
			if _, ok := sentimentCounts[analysis.Score]; !ok {
				sentimentCounts[analysis.Score] = 1
				continue
			}
			sentimentCounts[analysis.Score]++
		}

		// Form the output.
		metrics := Metrics{
			Comments: CommentMetrics{
				TotalCount:    len(thread.Comments),
				PositiveCount: sentimentCounts[1],
				NegativeCount: sentimentCounts[0],
			},
			Users: UserMetrics{
				TotalCount: len(thread.Users),
			},
			Actions: ActionMetrics{
				TotalCount: len(thread.Actions),
			},
		}

		// Count the likes and flags.
		var likes int
		var flags int
		for _, action := range thread.Actions {

			switch action.Type {
			case "like":
				likes += action.Count
			case "flag":
				flags += action.Count
			}
		}

		// Add the remaining action stats.
		metrics.Actions.LikeCount = likes
		metrics.Actions.FlagCount = flags

		// Marshal the output.
		outputData, err := json.MarshalIndent(metrics, "", "    ")
		if err != nil {
			return err
		}

		// Save the output data to a JSON file.
		if err := ioutil.WriteFile(filepath.Join(outDir, info.Name()), outputData, 0644); err != nil {
			return err
		}

		return nil

	}); err != nil {
		log.Fatal(err)
	}
}
