package main

import (
	"github.com/google/uuid"
)

type Item struct {
	ID     uuid.UUID
	name   string
	scores []int64
}

// func newItem() *Item {
// 	return &Item{
// 		ID: uuid
// 	}
// }
