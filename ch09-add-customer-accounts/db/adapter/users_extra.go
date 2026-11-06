package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"depot/app/models"
)

// FindUserByEmail finds a customer by their normalized email.
func (s *Store) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	row, err := s.q.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, models.ErrNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("store: find user by email: %w", err)
	}
	return userFromRow(row), nil
}
