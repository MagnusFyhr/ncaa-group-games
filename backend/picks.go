// session_store.go
package main

import "sync"

type Picks struct {
	VenmoName string  `json:"venmoName"`
	TeamName  string  `json:"teamName"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Pin       string  `json:"pin"`
	Picks     []*Team `json:"picks"`
}

type PicksStore struct {
	sync.Mutex
	Picks map[string]*Picks
}

func (store *PicksStore) CreatePicks(
	venmoName string,
	teamName string,
	firstName string,
	lastName string,
	pin string,
	picks []*Team,
) string {
	store.Lock()
	defer store.Unlock()

	store.Picks[venmoName] = &Picks{
		VenmoName: venmoName,
		TeamName:  teamName,
		FirstName: firstName,
		LastName:  lastName,
		Pin:       pin,
		Picks:     picks}

	return venmoName
}

func (store *PicksStore) GetAllPicks() []*Picks {
	store.Lock()
	defer store.Unlock()

	picks := make([]*Picks, 0, len(store.Picks))
	for _, pick := range store.Picks {
		picks = append(picks, pick)
	}

	return picks
}

func (store *PicksStore) GetPicks(venmoName string) (*Picks, bool) {
	store.Lock()
	defer store.Unlock()

	picksObj, exists := store.Picks[venmoName]
	return picksObj, exists
}

func (store *PicksStore) GetPicksWithPin(venmoName string, pin string) (*Picks, bool) {
	store.Lock()
	defer store.Unlock()

	picksObj, exists := store.Picks[venmoName]
	if exists {
		if picksObj.Pin != pin {
			return nil, false
		}
	}
	return picksObj, exists
}

func (store *PicksStore) UpdatePicks(
	venmoName string,
	teamName string,
	pin string,
	picks []*Team,
) (*Picks, bool) {
	store.Lock()
	defer store.Unlock()

	picksObj, exists := store.Picks[venmoName]

	if !exists {
		return nil, exists
	}

	if picksObj.Pin != pin {
		return nil, exists
	}

	if teamName != "" {
		store.Picks[venmoName].TeamName = teamName
	}
	store.Picks[venmoName].Picks = picks
	picksObj, exists = store.Picks[venmoName]

	return picksObj, exists
}

func (store *PicksStore) DeletePicks(
	venmoName string,
	pin string,
) (bool, bool) {
	store.Lock()
	defer store.Unlock()

	picksObj, exists := store.Picks[venmoName]

	if !exists {
		return false, exists
	}

	if picksObj.Pin != pin {
		return false, exists
	}

	// delete
	delete(store.Picks, venmoName)

	return true, exists
}

func (store *PicksStore) DeleteAllPicks() bool {
	store.Lock()
	defer store.Unlock()

	// Reassign an empty map
	store.Picks = make(map[string]*Picks)

	return true
}
