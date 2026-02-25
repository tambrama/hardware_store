package category

import (
	"github.com/google/uuid"
)

type Category struct {
	CategoryID uuid.UUID
	Category   string
}
