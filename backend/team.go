// questions.go
package main

import (
	"encoding/json"
	"os"
	"sync"
)

type Team struct {
	UID        int    `json:"id"`
	TeamName   string `json:"teamName"`
	Seed       int    `json:"seed"`
	IsGameOver bool   `json:"isGameOver"`
	Won        int    `json:"won"`
}

type TeamsStore struct {
	sync.Mutex
	Teams []*Team
}

func loadTeams() ([]*Team, error) {
	fileBytes, err := os.ReadFile("teams.json")
	if err != nil {
		return nil, err
	}

	var teams []Team
	if err := json.Unmarshal(fileBytes, &teams); err != nil {
		return nil, err
	}

	// Convert each element of teams to a pointer to Team
	var teamsPtrs []*Team
	for i := range teams {
		teamsPtrs = append(teamsPtrs, &teams[i])
	}

	return teamsPtrs, nil
}

func (store *TeamsStore) GetTeams() []*Team {
	store.Lock()
	defer store.Unlock()

	teams := make([]*Team, 0, len(store.Teams))
	for _, team := range store.Teams {
		teams = append(teams, team)
	}

	return teams
}

func (store *TeamsStore) UpdateTeam(
	teamName string,
	isGameOver bool,
	won int,
) (*Team, bool) {
	store.Lock()
	defer store.Unlock()

	team, exists := store.GetTeam(teamName)
	if !exists {
		return team, exists
	}

	team.IsGameOver = isGameOver
	team.Won = won

	return team, true
}

func (store *TeamsStore) GetTeam(teamName string) (*Team, bool) {
	store.Lock()
	defer store.Unlock()

	for _, team := range store.Teams {
		if team.TeamName == teamName {
			return team, true
		}
	}

	return nil, false
}
