package models

import (
	"html"
	"time"

	"github.com/satori/uuid"
)

type PVZ struct {
	Id               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city"`
}

func (p *PVZ) Sanitize() {
	p.City = html.EscapeString(p.City)
}
