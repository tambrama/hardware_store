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

func (s *storage) SaveSession(ctx context.Context, session model.Session) error {
	const op = "storage.sqlite.SaveSession"
	query := `INSERT INTO refresh_tokens (session_id, user_id, app_id, token_hash, expires_at) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, session.SessionID, session.UserID, session.AppID, session.HashToken, session.ExpiresAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("%s: %w", op, model.ErrRefreshTokenExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *storage) GetSession(ctx context.Context, tokenHash string) (model.Session, error) {
	const op = "storage.sqlite.GetSession"
	var session model.Session
	query := `SELECT session_id, user_id, app_id, token_hash, expires_at FROM refresh_tokens WHERE token_hash = ? AND expires_at > ?`
	row := s.db.QueryRowContext(ctx, query, tokenHash, time.Now())
	err := row.Scan(&session.SessionID, &session.UserID, &session.AppID, &session.HashToken, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Session{}, fmt.Errorf("%s: %w", op, model.ErrRefreshTokenNotFound)
		}
		return model.Session{}, fmt.Errorf("%s: %w", op, err)
	}
	return session, nil
}

func (s *storage) DeleteSession(ctx context.Context, userID uuid.UUID, appID string) error {
	const op = "storage.sqlite.DeleteSession"
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

func (s *storage) UpdateSession(ctx context.Context, oldToken, newToken string, duration time.Time) error {
	const op = "storage.sqlite.UpdateSession"
	query := `UPDATE refresh_tokens SET token_hash = ?, expires_at = ? 
	WHERE token_hash = ? AND expires_at > ?`
	res, err := s.db.ExecContext(ctx, query, newToken, duration, oldToken, time.Now())
	if err != nil {
		return fmt.Errorf("%s: exec update: %w", op, err)
	}
	row, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: check rows affected: %w", op, err)
	}
	if row == 0 {
		return model.ErrRefreshTokenNotFound
	}
	return nil
}
