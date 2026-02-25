package client

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ClientID         uuid.UUID `validate:"required"`
	Name             string    `validate:"required,min=2,max=50"`
	Surname          string    `validate:"required,min=2,max=50"`
	Birthday         time.Time `validate:"datetime=2006-01-02"`
	Gender           string    `validate:"oneof=male female"`
	RegistrationDate time.Time
	AddressID        uuid.UUID
}
