package models

import (
	"time"

	"github.com/satori/uuid"
)

type Product struct {
	Id            uuid.UUID `json:"id"`
	ReceptionDate time.Time `json:"reception_date"`
	Category      string    `json:"category"`
	ReceptionId   uuid.UUID `json:"reception_id"`
}

type Reception struct {
	Id            uuid.UUID `json:"id"`
	ReceptionDate time.Time `json:"reception_date"`
	PvzId         uuid.UUID `json:"pvzid"`
	Products      []Product `json:"products"`
	Status        string    `json:"status"`
}
