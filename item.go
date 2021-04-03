package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Item struct {
	ID     uuid.UUID
	name   string
	scores []int64
}

func newItem(n string, l *zap.Logger) *Item {
	ID, err := uuid.NewRandom()
	if err != nil {
		l.Error("Error while creating new item id: %s", zap.Error(err))
	}
	return &Item{
		ID:     ID,
		name:   n,
		scores: nil,
	}
}
