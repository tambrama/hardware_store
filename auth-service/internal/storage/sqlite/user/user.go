package userstorage

import (
	"auth-service/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type storage struct {
	db *sql.DB
}	
func NewStorage(db *sql.DB) *storage {
	return &storage{db: db}
}
func (s *storage) SaveUser(ctx context.Context, user model.Users) error {
	const op = "storage.sqlite.SaveUser"
	query := `INSERT INTO users (id, mail, hash_password, name, surname, phone_number) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query, user.ID.String(), user.Mail, user.Password, user.Name, user.Surname, user.PhoneNumber)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("%s: %w", op, model.ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (s *storage) GetUserByEmail(ctx context.Context, email string) (model.Users, error) {
	const op = "storage.sqlite.GetUserByEmail"
	var user model.Users
	query := `SELECT id, mail, hash_password, name, surname, phone_number FROM users WHERE mail = ?`
	row := s.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&user.ID, &user.Mail, &user.Password, &user.Name, &user.Surname, &user.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}
	return user, err
}
func (s *storage) UpdateUserPassword(ctx context.Context, email string, newPassword []byte) error {
	const op = "storage.sqlite.UpdateUserPassword"
	query := `UPDATE users SET hash_password = ? WHERE mail = ?`
	result, err := s.db.ExecContext(ctx, query, newPassword, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: user with email %s not found", op, email)
	}
	return nil
}
