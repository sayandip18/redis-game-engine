package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	instance *redis.Client
	once sync.Once
)

type Config struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// Get returns the singleton Redis client, initializing it on first call.
func Get(cfg Config) *redis.Client {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr: cfg.Addr,
			Password: cfg.Password,
			DB: cfg.DB,
			DialTimeout: cfg.DialTimeout,
			ReadTimeout: cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,

			// Connection pool tuning — important for concurrent modules
			PoolSize: 10,
			MinIdleConns: 3,
			MaxRetries: 3,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatalf("❌ Redis connection failed: %v", err)
		}

		info := rdb.Info(ctx, "server").Val()
		fmt.Println("✅ Redis connected successfully")
		fmt.Println(info)

		instance = rdb
	})
	return instance
}


// Close gracefully closes the connection pool.
func Close() error {
	if instance != nil {
		return instance.Close()
	}
	return nil
}