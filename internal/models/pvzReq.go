package models

import (
	"html"

	"github.com/satori/uuid"
)

// easyjson:json
type AddProductReq struct {
	Type  string    `json:"type"`
	PvzId uuid.UUID `json:"pvzId"`
}

func (p *AddProductReq) Sanitize() {
	p.Type = html.EscapeString(p.Type)
}
