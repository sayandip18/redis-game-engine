package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/sayandip18/redis-game-engine/internal/client"
)

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// CreateSession stores a userID→token mapping in Redis with the given TTL and returns the token.
func CreateSession(userID string, ttl string) (string, error) {
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		return "", fmt.Errorf("invalid ttl %q: %w", ttl, err)
	}

	token, err := generateToken()
	if err != nil {
		return "", err
	}

	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	if err := rdb.Set(ctx, "session:"+token, userID, duration).Err(); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return token, nil
}

// ValidateSession looks up the token in Redis and returns the associated userID.
func ValidateSession(token string) (string, error) {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	userID, err := rdb.Get(ctx, "session:"+token).Result()
	if err != nil {
		return "", fmt.Errorf("session not found or expired: %w", err)
	}

	return userID, nil
}

// RevokeSession deletes the session token from Redis.
func RevokeSession(token string) error {
	rdb := client.Get(client.DefaultConfig())
	ctx := context.Background()

	n, err := rdb.Del(ctx, "session:"+token).Result()
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}
