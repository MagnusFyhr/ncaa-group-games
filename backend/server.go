// game_server.go
package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	PointsPerSeedStore *PointsPerSeedStore
	TeamsStore         *TeamsStore
	PicksStore         *PicksStore
	LeaderboardStore   *LeaderboardStore
	PicksDeadline      time.Time
	Mutex              sync.Mutex
}

func NewServer(picksDeadline time.Time, pointsPerSeedData map[int]int, teamsData []*Team) *Server {
	return &Server{
		TeamsStore: &TeamsStore{
			Teams: teamsData,
		},
		PicksStore: &PicksStore{
			Picks: make(map[string]*Picks),
		},
		PointsPerSeedStore: &PointsPerSeedStore{
			PointsPerSeed: pointsPerSeedData,
		},
		LeaderboardStore: &LeaderboardStore{
			Leaderboard:   make([]*Leaderboard, 0),
			PointsPerSeed: pointsPerSeedData,
		},
		PicksDeadline: picksDeadline,
	}
}

func (s *Server) GetPicksDeadline(c *gin.Context) {

	picksDeadline := s.PicksDeadline.UTC().Format(time.RFC3339)

	c.JSON(http.StatusOK, gin.H{
		"picksDeadline": picksDeadline,
	})
}

func (s *Server) GetNumberOfPlayers(c *gin.Context) {

	numberOfPlayers := len(s.PicksStore.Picks)

	c.JSON(http.StatusOK, gin.H{
		"numberOfPlayers": numberOfPlayers,
	})
}

func (s *Server) GetTeams(c *gin.Context) {
	teams := s.TeamsStore.GetTeams()
	c.JSON(http.StatusOK, gin.H{
		"teams": teams,
	})
}

func (s *Server) GetPointsPerSeed(c *gin.Context) {
	pointsPerSeedData := s.PointsPerSeedStore.GetPointsPerSeeds()
	c.JSON(http.StatusOK, gin.H{
		"pointsPerSeed": pointsPerSeedData,
	})
}

func (s *Server) UpdateTeam(c *gin.Context) {
	var updateTeamPayload struct {
		TeamName   string `json:"teamName"`
		IsGameOver bool   `json:"isGameOver"`
		Won        int    `json:"won"`
	}
	if err := c.ShouldBindJSON(&updateTeamPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	_, exists := s.TeamsStore.GetTeam(updateTeamPayload.TeamName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get Team"})
		return
	}

	if time.Now().After(s.PicksDeadline) {
		c.JSON(http.StatusOK, gin.H{"error": "Failed, deadline has passed"})
		return
	}

	s.Mutex.Lock()
	s.TeamsStore.UpdateTeam(updateTeamPayload.TeamName, updateTeamPayload.IsGameOver, updateTeamPayload.Won)
	s.Mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (s *Server) SetPicks(c *gin.Context) {
	venmoName := c.Param("venmoName")
	var createPicksPayload struct {
		Pin       string   `json:"pin"`
		TeamName  string   `json:"teamName"`
		FirstName string   `json:"firstName"`
		LastName  string   `json:"lastName"`
		Picks     []string `json:"picks"`
	}
	if err := c.ShouldBindJSON(&createPicksPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if venmoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid venmo"})
		return
	}

	if createPicksPayload.Pin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pin"})
		return
	}

	if createPicksPayload.TeamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teamName"})
		return
	}
	if createPicksPayload.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid firstName"})
		return
	}
	if createPicksPayload.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lastName"})
		return
	}
	if time.Now().After(s.PicksDeadline) {
		c.JSON(http.StatusOK, gin.H{"error": "Failed, deadline has passed"})
		return
	}

	// convert picks []string to picks []*Team
	var picksPtrs []*Team
	for i := 0; i < len(createPicksPayload.Picks); i++ {
		teamName := createPicksPayload.Picks[i]
		team, exists := s.TeamsStore.GetTeam(teamName)
		// is it a valid team name
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(`Invalid team name: %s`, teamName)})
			return
		}
		// is it valid seed
		if team.Seed-1 != i {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is a %d seed; input as %d seed", teamName, team.Seed, i)})
			return
		}
		// append to list
		picksPtrs = append(picksPtrs, team)
	}

	// check if user has already made picks
	picks, exists := s.PicksStore.GetPicks(venmoName)
	s.Mutex.Lock()
	if !exists {
		s.PicksStore.CreatePicks(
			venmoName,
			createPicksPayload.TeamName,
			createPicksPayload.FirstName,
			createPicksPayload.LastName,
			createPicksPayload.Pin,
			picksPtrs)
	} else {
		response, _ := s.PicksStore.UpdatePicks(
			venmoName,
			createPicksPayload.TeamName,
			createPicksPayload.Pin,
			picksPtrs)
		if response == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update Pick; Invalid PIN"})
			return
		}
	}
	s.Mutex.Unlock()

	picks, _ = s.PicksStore.GetPicks(venmoName)
	var picksStrings []string
	for _, team := range picksPtrs {
		picksStrings = append(picksStrings, team.TeamName)
	}
	c.JSON(http.StatusOK, gin.H{
		"venmoName": picks.VenmoName,
		"firstName": picks.FirstName,
		"lastName":  picks.LastName,
		"teamName":  picks.TeamName,
		"picks":     picksStrings,
	})
}

func (s *Server) GetPick(c *gin.Context) {
	venmoName := c.Param("venmoName")
	pin := c.Param("pin")

	if pin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pin"})
		return
	}
	if venmoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid venmo name"})
		return
	}

	picks, exists := s.PicksStore.GetPicksWithPin(venmoName, pin)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get Pick"})
		return
	}

	// convert picks []*Team to picks []string
	var picksStrings []string
	for _, team := range picks.Picks {
		picksStrings = append(picksStrings, team.TeamName)
	}

	c.JSON(http.StatusOK, gin.H{
		"venmoName": picks.VenmoName,
		"firstName": picks.FirstName,
		"lastName":  picks.LastName,
		"teamName":  picks.TeamName,
		"picks":     picksStrings,
	})
}

func (s *Server) DeletePick(c *gin.Context) {
	venmoName := c.Param("venmoName")
	var getPicksPayload struct {
		Pin string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&getPicksPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Delete the pick from memory
	deleted, exists := s.PicksStore.DeletePicks(venmoName, getPicksPayload.Pin)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete Pick; Invalid Venmo"})
		return
	}
	if exists && !deleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete Pick; Invalid PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (s *Server) ClearData(c *gin.Context) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	success := s.PicksStore.DeleteAllPicks()
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete all picks"})
		return
	}
	success = s.LeaderboardStore.DeleteAllEntries()
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to clear leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (s *Server) GetLeaderboard(c *gin.Context) {
	if time.Now().Before(s.PicksDeadline) {
		c.JSON(http.StatusOK, gin.H{"error": "Failed, deadline has not passed"})
		return
	}

	// refresh leaderboard
	picks := s.PicksStore.GetAllPicks()
	success := s.LeaderboardStore.refresh(picks)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to refresh leaderboard"})
		return
	}

	// get leaderboard
	leaderboard := s.LeaderboardStore.GetLeaderboard()

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
	})
}

func (s *Server) UpdateServer(c *gin.Context) {
	var updateServerPayload struct {
		PicksDeadline string `json:"picksDeadline"`
	}
	if err := c.ShouldBindJSON(&updateServerPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Parse the ISO 8601 date-time string into a time.Time object
	picksDeadline, err := time.Parse(time.RFC3339, updateServerPayload.PicksDeadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Error parsing time: " + err.Error(),
		})
		return
	}

	s.Mutex.Lock()
	s.PicksDeadline = picksDeadline
	s.Mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
