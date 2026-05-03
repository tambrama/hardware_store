package photo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hardware_store/internal/model/images"
	"io"
	"log/slog"
	"net/http"
)

type PhotoServiceClient struct {
	URL        string
	httpClient *http.Client
	logger     *slog.Logger
}

func NewPhotoServiceClient(url string, logger *slog.Logger) *PhotoServiceClient {
	return &PhotoServiceClient{
		URL:        url,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

func (c *PhotoServiceClient) CreatePhoto(ctx context.Context, photoData []byte) (images.ImageResponse, error) {
	const op = "PhotoServiceClient.CreatePhoto"
	req, err := http.NewRequestWithContext(ctx, "POST", c.URL+"/images", bytes.NewReader(photoData))
	if err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to create request", slog.Any("error", err))
		return images.ImageResponse{}, fmt.Errorf("%s: create request: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to send request", slog.Any("error", err))
		return images.ImageResponse{}, fmt.Errorf("%s: send request: %w", op, err)
	}
	defer resp.Body.Close()

	var img images.ImageResponse

	if err := json.NewDecoder(resp.Body).Decode(&img); err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to decode response", slog.Any("error", err))
		return images.ImageResponse{}, fmt.Errorf("%s: decode response: %w", op, err)
	}

	return img, nil
}

func (c *PhotoServiceClient) UpdatePhoto(ctx context.Context, photoData []byte, imgID string) (images.ImageResponse, error) {
	const op = "PhotoServiceClient.UpdatePhoto"

	req, err := http.NewRequestWithContext(ctx, "PUT", c.URL+"/images/"+imgID, bytes.NewReader(photoData))
	if err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to create request", slog.Any("error", err))
		return images.ImageResponse{}, fmt.Errorf("%s: create request: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to send request", slog.Any("error", err))
		return images.ImageResponse{}, fmt.Errorf("%s: send request: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return images.ImageResponse{}, fmt.Errorf("%s: unexpected status %d: %s", op, resp.StatusCode, string(body))
	}

	var img images.ImageResponse

	if err := json.NewDecoder(resp.Body).Decode(&img); err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to decode response", slog.Any("error", err))
		return images.ImageResponse{}, fmt.Errorf("%s: decode response: %w", op, err)
	}

	return img, nil
}

func (c *PhotoServiceClient) GetPhoto(ctx context.Context, imgID string) ([]byte, error) {
	const op = "PhotoServiceClient.GetPhoto"

	req, err := http.NewRequestWithContext(ctx, "GET", c.URL+"/images/"+imgID, nil)
	if err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to create request", slog.Any("error", err))
		return nil, fmt.Errorf("%s: create request: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to send request", slog.Any("error", err))
		return nil, fmt.Errorf("%s: send request: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s: unexpected status %d: %s", op, resp.StatusCode, string(body))
	}

	photoData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: read response body: %w", op, err)
	}

	return photoData, nil
}

func (c *PhotoServiceClient) DeletePhoto(ctx context.Context, imgID string) error {
	const op = "PhotoServiceClient.DeletePhoto"

	req, err := http.NewRequestWithContext(ctx, "DELETE", c.URL+"/images/"+imgID, nil)
	if err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to create request", slog.Any("error", err))
		return fmt.Errorf("%s: create request: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Log(context.Background(), slog.LevelError, "Failed to send request", slog.Any("error", err))
		return fmt.Errorf("%s: send request: %w", op, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s: unexpected status %d: %s", op, resp.StatusCode, string(body))
	}
	return nil
}
