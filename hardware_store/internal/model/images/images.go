package images

import (
	"github.com/google/uuid"
)

type Images struct {
	ImageID uuid.UUID
	Image   []byte
}
