package player

import (
	"context"
	"fmt"

	"github.com/sayandip18/redis-game-engine/internal/client"
)

type LeaderboardEntry struct {
	Rank     int
	PlayerID string
	MMR      int64
}

// GetRank returns the 1-indexed global rank of a player (highest MMR = rank 1).
func GetRank(playerID string) (int64, error) {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	rank, err := rdb.ZRevRank(ctx, leaderboardKey, playerID).Result()
	if err != nil {
		return 0, fmt.Errorf("get rank: %w", err)
	}
	return rank + 1, nil
}

// GetGlobalTop returns the top n players by MMR from the global leaderboard.
func GetGlobalTop(n int) ([]LeaderboardEntry, error) {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	results, err := rdb.ZRevRangeWithScores(ctx, leaderboardKey, 0, int64(n-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("get global top: %w", err)
	}

	entries := make([]LeaderboardEntry, len(results))
	for i, z := range results {
		entries[i] = LeaderboardEntry{
			Rank:     i + 1,
			PlayerID: z.Member.(string),
			MMR:      int64(z.Score),
		}
	}
	return entries, nil
}
