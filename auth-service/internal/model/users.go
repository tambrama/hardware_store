package model

import "github.com/google/uuid"

type Users struct {
	ID          uuid.UUID
	Name        string
	Surname     string
	Mail        string
	PhoneNumber string
	Password    []byte
}
