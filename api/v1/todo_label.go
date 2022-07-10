package v1

import (
	"github.com/gofrs/uuid"
)

type TodoLabel struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type TodoLabels []TodoLabel
