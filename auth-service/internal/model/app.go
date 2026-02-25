package model

import "github.com/google/uuid"

type App struct {
	ID          uuid.UUID
	Name        string
}