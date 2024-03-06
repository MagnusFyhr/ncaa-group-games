// questions.go
package main

import (
	"sync"
)

type Leaderboard struct {
	Id        int      `json:"id"`
	TeamName  string   `json:"teamName"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Picks     []string `json:"picks"`
	Points    int      `json:"points"`
	MaxPoints int      `json:"maxPoints"`
}

type LeaderboardStore struct {
	sync.Mutex
	Leaderboard   []*Leaderboard
	PointsPerSeed map[int]int
}

func (store *LeaderboardStore) refresh(picksArr []*Picks) bool {
	store.Lock()
	defer store.Unlock()

	// create leaderboard entries
	store.Leaderboard = []*Leaderboard{} // Reassign an empty slice
	for id, picks := range picksArr {
		store.Leaderboard = append(store.Leaderboard, buildLeaderboardEntry(id+1, picks, store.PointsPerSeed))
	}

	return true
}

func (store *LeaderboardStore) GetLeaderboard() []*Leaderboard {
	store.Lock()
	defer store.Unlock()

	return store.Leaderboard
}

func buildLeaderboardEntry(id int, picks *Picks, pointsPerSeedData map[int]int) *Leaderboard {
	currPoints := 0
	maxPoints := 0
	for _, team := range picks.Picks {
		currPoints += team.Won * pointsPerSeedData[team.Seed]
		maxPoints += team.Won * pointsPerSeedData[team.Seed]
		if !team.IsGameOver {
			maxPoints += pointsPerSeedData[team.Seed]
		}
	}

	// convert picks []*Team to picks []string
	var picksStrings []string
	for _, team := range picks.Picks {
		picksStrings = append(picksStrings, team.TeamName)
	}

	return &Leaderboard{
		Id:        id,
		TeamName:  picks.TeamName,
		FirstName: picks.FirstName,
		LastName:  picks.LastName,
		Picks:     picksStrings,
		Points:    currPoints,
		MaxPoints: maxPoints,
	}
}

func (store *LeaderboardStore) DeleteAllEntries() bool {
	store.Lock()
	defer store.Unlock()

	// Reassign an empty map
	store.Leaderboard = make([]*Leaderboard, 0)

	return true
}
