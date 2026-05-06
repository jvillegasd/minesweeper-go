package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const maxScoresPerDifficulty = 5

type ScoreEntry struct {
	Seconds int       `json:"seconds"`
	Date    time.Time `json:"date"`
}

type Leaderboard struct {
	Beginner     []ScoreEntry `json:"beginner"`
	Intermediate []ScoreEntry `json:"intermediate"`
	Expert       []ScoreEntry `json:"expert"`
}

var leaderboard = &Leaderboard{}

func leaderboardPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".minesweeper-go-scores.json"
	}
	return filepath.Join(home, ".minesweeper-go-scores.json")
}

func loadLeaderboard() *Leaderboard {
	lb := &Leaderboard{}
	data, err := os.ReadFile(leaderboardPath())
	if err != nil {
		return lb
	}
	_ = json.Unmarshal(data, lb)
	return lb
}

func saveLeaderboard(lb *Leaderboard) {
	data, err := json.MarshalIndent(lb, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(leaderboardPath(), data, 0o644)
}

func (lb *Leaderboard) entries(d Difficulty) []ScoreEntry {
	switch d {
	case Beginner:
		return lb.Beginner
	case Intermediate:
		return lb.Intermediate
	case Expert:
		return lb.Expert
	}
	return nil
}

func (lb *Leaderboard) setEntries(d Difficulty, entries []ScoreEntry) {
	switch d {
	case Beginner:
		lb.Beginner = entries
	case Intermediate:
		lb.Intermediate = entries
	case Expert:
		lb.Expert = entries
	}
}

func (lb *Leaderboard) record(d Difficulty, seconds int) (rank int, isBest bool) {
	entry := ScoreEntry{Seconds: seconds, Date: time.Now()}
	entries := append(lb.entries(d), entry)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Seconds < entries[j].Seconds
	})
	if len(entries) > maxScoresPerDifficulty {
		entries = entries[:maxScoresPerDifficulty]
	}
	lb.setEntries(d, entries)

	rank = 0
	for i, e := range entries {
		if e.Seconds == seconds && e.Date.Equal(entry.Date) {
			rank = i + 1
			break
		}
	}
	isBest = rank == 1
	return rank, isBest
}

func (lb *Leaderboard) best(d Difficulty) (int, bool) {
	entries := lb.entries(d)
	if len(entries) == 0 {
		return 0, false
	}
	return entries[0].Seconds, true
}
