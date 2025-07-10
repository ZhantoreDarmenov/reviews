package models

import (
	"time"
)

type Reviews struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Photo       string     `json:"photo"`
	Description string     `json:"description"`
	Rating      int        `json:"rating"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
