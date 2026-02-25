package jwtstorage

import (
	"auth-service/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *storage {
	return &storage{db: db}
}

func (s *storage) SaveRefreshToken(ctx context.Context, userID uuid.UUID, appID string, refreshToken string, expiresAt time.Time) error {
	const op = "storage.sqlite.SaveRefreshToken"
	query := `INSERT OR REPLACE INTO refresh_tokens (token_hash, user_id, app_id, expires_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, refreshToken, userID, appID, expiresAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("%s: %w", op, model.ErrRefreshTokenExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *storage) GetRefreshToken(ctx context.Context, userID uuid.UUID, appID string) (string, error) {
	const op = "storage.sqlite.GetRefreshToken"
	var refreshToken string
	query := `SELECT refresh_token FROM refresh_tokens WHERE user_id = ? AND app_id = ?`
	row := s.db.QueryRowContext(ctx, query, userID, appID)
	err := row.Scan(&refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, model.ErrRefreshTokenNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return refreshToken, nil
}

func (s *storage) DeleteRefreshToken(ctx context.Context, userID uuid.UUID, appID string) error {
	const op = "storage.sqlite.DeleteRefreshToken"
	query := `DELETE FROM refresh_tokens WHERE user_id = ? AND app_id = ?`
	result, err := s.db.ExecContext(ctx, query, userID, appID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: no refresh token found for user ID %s and app ID %s", op, userID.String(), appID)
	}
	return nil
}
