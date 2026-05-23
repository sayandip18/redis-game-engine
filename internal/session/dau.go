package session

import (
	"context"
	"fmt"

	"github.com/sayandip18/redis-game-engine/internal/client"
)

// MarkLogin records that userID was active on the given date (format: "2006-01-02").
// Idempotent — calling it multiple times for the same user and date is safe.
func MarkLogin(userID string, date string) error {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	if err := rdb.SAdd(ctx, "dau:"+date, userID).Err(); err != nil {
		return fmt.Errorf("mark login: %w", err)
	}

	return nil
}

// GetDAU returns the count of unique active users on the given date (format: "2006-01-02").
func GetDAU(date string) (int64, error) {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	count, err := rdb.SCard(ctx, "dau:"+date).Result()
	if err != nil {
		return 0, fmt.Errorf("get dau: %w", err)
	}

	return count, nil
}
