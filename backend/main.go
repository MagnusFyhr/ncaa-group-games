// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Setup the server
	router, err := setupServer()
	if err != nil {
		log.Fatalf("Server setup failed: %v", err)
	}

	// set port to PORT or 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Println("Server starting on port " + port)
	log.Fatal(router.Run(":" + port))
}

func setupServer() (*gin.Engine, error) {

	pointsPerSeedData, err := loadPointsPerSeed()
	if err != nil {
		return nil, err
	}
	teamsData, err := loadTeams()
	if err != nil {
		return nil, err
	}

	// Parse the ISO 8601 date-time string into a time.Time object
	picksDeadline, err := time.Parse(time.RFC3339, "2024-03-21T16:00:00.000Z")
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return nil, nil
	}
	// this is creating a game server
	server := NewServer(picksDeadline, pointsPerSeedData, teamsData)

	// Create Gin router and setup routes
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	config := cors.DefaultConfig()

	// allow all origins
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/teams", server.GetTeams)
	router.GET("/points-per-seed", server.GetPointsPerSeed)
	router.GET("/deadline", server.GetPicksDeadline)
	router.GET("/number-of-players", server.GetNumberOfPlayers)

	router.PUT("/picks/:venmoName", server.SetPicks)
	router.GET("/picks/:venmoName/:pin", server.GetPick)
	router.DELETE("/picks/:venmoName", server.DeletePick)

	router.GET("/leaderboard", server.GetLeaderboard)
	router.PATCH("/teams", server.UpdateTeam)

	router.PATCH("/server", server.UpdateServer)

	return router, nil
}
