package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"hardware_store/internal/model/images"
	"hardware_store/internal/redis"
	"time"

	"github.com/google/uuid"
)

type ImageCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewImageCache(client *redis.Client, prefix string, ttl time.Duration) *ImageCache {
	return &ImageCache{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}
}

func (c *ImageCache) Get(ctx context.Context, productID uuid.UUID) (images.Images, error) {
	key := fmt.Sprintf("%s%s", c.prefix, productID)
	data, err := c.client.Get(ctx, key).Bytes()

	if err != nil {
		if err.Error() == "redis: nil" {
			return images.Images{}, nil
		}
		return images.Images{}, err
	}
	var img images.Images
	if err := json.Unmarshal(data, &img); err != nil {
		return images.Images{}, err
	}
	return img, nil
}

func (c *ImageCache) Set(ctx context.Context, productID uuid.UUID, images images.Images) error {
	key := fmt.Sprintf("%s%s", c.prefix, productID)
	data, err := json.Marshal(images)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *ImageCache) Delete(ctx context.Context, productID uuid.UUID) error {
	key := fmt.Sprintf("%s%s", c.prefix, productID)
	return c.client.Del(ctx, key).Err()
}

func (c *ImageCache) Clear(ctx context.Context) error {
	iter := c.client.Scan(ctx, 0, fmt.Sprintf("%s*", c.prefix), 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
