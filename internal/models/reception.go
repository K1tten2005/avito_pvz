package models

import (
	"time"

	"github.com/satori/uuid"
)

// easyjson:json
type Product struct {
	Id            uuid.UUID `json:"id"`
	ReceptionTime time.Time `json:"reception_time"`
	Category      string    `json:"category"`
	ReceptionId   uuid.UUID `json:"reception_id"`
}

// easyjson:json
type Reception struct {
	Id            uuid.UUID `json:"id"`
	ReceptionTime time.Time `json:"reception_time"`
	PvzId         uuid.UUID `json:"pvzid"`
	Products      []Product `json:"products"`
	Status        string    `json:"status"`
}
