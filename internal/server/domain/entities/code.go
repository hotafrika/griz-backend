package entities

import "time"

type Code struct {
	ID        uint64
	UserID    uint64
	SrcURL    string
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
