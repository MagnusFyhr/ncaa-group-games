// questions.go
package main

import (
	"encoding/json"
	"os"
	"sync"
)

type PointsPerSeedStore struct {
	sync.Mutex
	PointsPerSeed map[int]int
}

func loadPointsPerSeed() (map[int]int, error) {
	fileBytes, err := os.ReadFile("points.json")
	if err != nil {
		return nil, err
	}

	var pointsPerSeed map[int]int
	if err := json.Unmarshal(fileBytes, &pointsPerSeed); err != nil {
		return nil, err
	}

	return pointsPerSeed, nil
}

func (store *PointsPerSeedStore) GetPointsPerSeeds() map[int]int {
	store.Lock()
	defer store.Unlock()

	return store.PointsPerSeed
}

func (store *PointsPerSeedStore) GetPointsPerSeed(seed int) (int, bool) {
	store.Lock()
	defer store.Unlock()

	points, exists := store.PointsPerSeed[seed]

	if !exists {
		return 0, exists
	}

	return points, exists
}
