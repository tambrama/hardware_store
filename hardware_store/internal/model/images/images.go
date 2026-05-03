package images

import (
	"github.com/google/uuid"
)

type Images struct {
	ImageID uuid.UUID `json:"id"`
	// Image   []byte
}

type ImageResponse struct {
	ImageID uuid.UUID `json:"id"`
	ShardID int `json:"shard,omitempty"`
}