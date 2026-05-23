package player

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sayandip18/redis-game-engine/internal/client"
)

const leaderboardKey = "leaderboard:global"

type Profile struct {
	UserID    string
	Username  string
	MMR       int64
	CreatedAt time.Time
	Stats     map[string]int64
}

func playerKey(userID string) string {
	return "player:" + userID
}

// CreatePlayer initialises a new player hash and adds them to the global leaderboard.
// If the player already exists, the existing profile is returned unchanged.
func CreatePlayer(userID, username string) (Profile, error) {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()
	key := playerKey(userID)

	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return Profile{}, fmt.Errorf("create player: %w", err)
	}
	if exists > 0 {
		return GetProfile(userID)
	}

	now := time.Now().UTC()

	pipe := rdb.Pipeline()
	pipe.HSet(ctx, key,
		"username", username,
		"mmr", 0,
		"created_at", now.Format(time.RFC3339),
	)
	pipe.ZAdd(ctx, leaderboardKey, redis.Z{Score: 0, Member: userID})
	if _, err := pipe.Exec(ctx); err != nil {
		return Profile{}, fmt.Errorf("create player: %w", err)
	}

	return Profile{
		UserID:    userID,
		Username:  username,
		MMR:       0,
		CreatedAt: now,
		Stats:     map[string]int64{},
	}, nil
}

// UpdateMMR adjusts a player's MMR by delta (positive or negative) and keeps
// the leaderboard sorted set in sync via a pipeline.
// Returns the new MMR value.
func UpdateMMR(userID string, delta int64) (int64, error) {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	pipe := rdb.Pipeline()
	newMMR := pipe.HIncrBy(ctx, playerKey(userID), "mmr", delta)
	pipe.ZIncrBy(ctx, leaderboardKey, float64(delta), userID)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, fmt.Errorf("update mmr: %w", err)
	}

	return newMMR.Val(), nil
}

// GetProfile fetches the full player profile, including all generic stats.
func GetProfile(userID string) (Profile, error) {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	data, err := rdb.HGetAll(ctx, playerKey(userID)).Result()
	if err != nil {
		return Profile{}, fmt.Errorf("get profile: %w", err)
	}
	if len(data) == 0 {
		return Profile{}, fmt.Errorf("player not found: %s", userID)
	}

	mmr, _ := strconv.ParseInt(data["mmr"], 10, 64)
	createdAt, _ := time.Parse(time.RFC3339, data["created_at"])

	stats := make(map[string]int64)
	for k, v := range data {
		if strings.HasPrefix(k, "stat:") {
			val, _ := strconv.ParseInt(v, 10, 64)
			stats[strings.TrimPrefix(k, "stat:")] = val
		}
	}

	return Profile{
		UserID:    userID,
		Username:  data["username"],
		MMR:       mmr,
		CreatedAt: createdAt,
		Stats:     stats,
	}, nil
}

// IncrementStat atomically increments a named stat for the player.
// The stat name is arbitrary (e.g. "wins", "kills", "matches_played").
func IncrementStat(userID, stat string, delta int64) error {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	if err := rdb.HIncrBy(ctx, playerKey(userID), "stat:"+stat, delta).Err(); err != nil {
		return fmt.Errorf("increment stat %q: %w", stat, err)
	}
	return nil
}
