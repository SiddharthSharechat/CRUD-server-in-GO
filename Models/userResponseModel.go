package Models

import (
	"time"
)

type UserResponse struct {
	Id        string
	Name      string
	Email     string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
