package appstorage

import (
	"auth-service/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *storage {
	return &storage{db: db}
}
func (s *storage) GetAppByID(ctx context.Context, appID string) (model.App, error) {
	const op = "storage.sqlite.GetAppByID"
	var app model.App
	query := `SELECT id, name FROM apps WHERE id = ?`
	row := s.db.QueryRowContext(ctx, query, appID)
	err := row.Scan(&app.ID, &app.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app, fmt.Errorf("%s: %w", op, model.ErrAppNotFound)
		}
		return app, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
