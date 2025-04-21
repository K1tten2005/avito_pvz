package models

import (
	"time"

	"github.com/satori/uuid"
)

// easyjson:json
type Product struct {
	Id          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionId uuid.UUID `json:"reception_id"`
}

// easyjson:json
type Reception struct {
	Id       uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzId    uuid.UUID `json:"pvzId"`
	Products []Product `json:"products"`
	Status   string    `json:"status"`
}

const(
	StatusInProgress = "in_progress"
	StatusClose = "close"
)
