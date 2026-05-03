package model

import "github.com/google/uuid"

type Users struct {
	ID          uuid.UUID `json:"id" db:"id"`
	GoogleID    *string `json:"google_id,omitempty" db:"google_id"`
	Name        string `json:"name" db:"name"`
	Surname     string `json:"surname" db:"surname"`
	Mail        string `json:"mail" db:"mail"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Password    *[]byte `json:"-" db:"hash_password"`
}

type GoogleUserInfo struct {
    GoogleID      string `json:"id"`
    Name        string `json:"given_name"`
	Surname     string `json:"family_name"`
	Mail        string `json:"email"`
	// PhoneNumber string `json:"phone_number"`
}