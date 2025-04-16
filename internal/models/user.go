package models

import (
	"html"

	"github.com/satori/uuid"
)

type User struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

func (u *User) Sanitize() {
	u.Email = html.EscapeString(u.Email)
}
