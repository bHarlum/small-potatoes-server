package main

import (
	"github.com/google/uuid"
)

type Room struct {
	ID    uuid.UUID
	items []Item
	owner string
	hub   *Hub
}

func newRoom(items []Item) *Room {
	ID, err := uuid.NewRandom()
	if err != nil {
		// log.Errorf("Error while creating new room id: %s", err)
	}
	return &Room{
		ID:    ID,
		items: items,
		owner: "",
		hub:   newHub(),
	}
}

func (r *Room) start() {
	// log.Infof("Starting room: %s", r.ID)
	r.hub.run()
}
