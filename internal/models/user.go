package models

import (
	"html"

	"github.com/satori/uuid"
)

// easyjson:json
type User struct {
	Id           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	PasswordHash []byte    `json:"-"`
}

const (
	RoleEmployee  = "employee"
	RoleModerator = "moderator"
)

func (u *User) Sanitize() {
	u.Email = html.EscapeString(u.Email)
}
