package models

import (
	"html"
	"time"

	"github.com/satori/uuid"
)

// easyjson:json
type PVZ struct {
	Id               uuid.UUID   `json:"id"`
	RegistrationDate time.Time   `json:"registrationDate"`
	City             string      `json:"city"`
	Receptions       []Reception `json:"receptions"`
}

func (p *PVZ) Sanitize() {
	p.City = html.EscapeString(p.City)
}
