package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Room struct {
	ID    uuid.UUID
	items []Item
	owner string
	hub   *Hub
}

func newRoom(i []string, o string, l *zap.Logger) *Room {
	ID, err := uuid.NewRandom()
	if err != nil {
		l.Error("Error while creating new room id", zap.Error(err))
	}
	items := make([]Item, len(i))
	for i, item := range i {
		items[i] = *newItem(item, l)
	}

	return &Room{
		ID:    ID,
		items: items,
		owner: o,
		hub:   newHub(),
	}
}

func (r *Room) Start() {
	r.hub.run()
}
