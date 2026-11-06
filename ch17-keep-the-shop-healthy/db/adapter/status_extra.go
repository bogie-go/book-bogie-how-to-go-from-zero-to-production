package adapter

import (
	"context"
	"fmt"

	"depot/app/models"
)

// QueueHealth counts the jobs that need a person's attention. River owns the
// river_job table, outside db/migrate, so sqlc can't see it; this is the one
// query the store runs by hand.
func (s *Store) QueueHealth(ctx context.Context) (models.QueueHealth, error) {
	var h models.QueueHealth
	err := s.pool.QueryRow(ctx, `
		SELECT
		  count(*) FILTER (WHERE state = 'available' AND scheduled_at < now() - interval '5 minutes'),
		  count(*) FILTER (WHERE state = 'discarded' AND finalized_at > now() - interval '1 day')
		FROM river_job`).Scan(&h.WaitingTooLong, &h.FailedToday)
	if err != nil {
		return models.QueueHealth{}, fmt.Errorf("store: queue health: %w", err)
	}
	return h, nil
}
