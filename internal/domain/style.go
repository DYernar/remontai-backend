// domain/style.go
package domain

import (
	"errors"
	"time"
)

// Style-related errors
var (
	ErrStyleNotFound = errors.New("style not found")
	ErrStyleExists   = errors.New("style already exists")
)

type StyleModel struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Image       string    `json:"image" db:"image"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
