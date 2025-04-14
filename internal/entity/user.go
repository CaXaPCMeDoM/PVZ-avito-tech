package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
}
