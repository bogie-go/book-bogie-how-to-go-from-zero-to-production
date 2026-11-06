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

// IsStaff reports whether the customer is on the shop's staff.
func (s *Store) IsStaff(ctx context.Context, userID int64) (bool, error) {
	staff, err := s.q.IsStaff(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("store: is staff: %w", err)
	}
	return staff, nil
}
