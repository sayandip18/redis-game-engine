package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sayandip18/redis-game-engine/internal/client"
)

func main() {
	cfg := client.DefaultConfig()
	rdb := client.Get(cfg)
	defer client.Close()

	ctx:= context.Background()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
		case "ping":
			result, err := rdb.Ping(ctx).Result()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		case "flush":
			fmt.Print("⚠️  Flushing ALL Redis data. Confirm? [yes/no]: ")
			var confirm string
			fmt.Scan(&confirm)
			if confirm == "yes" {
				rdb.FlushAll(ctx)
				fmt.Println("🔥 Database flushed.")
			}
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printHelp()

	}

}

func printHelp() {
	fmt.Println(`
Redis Game Engine CLI
─────────────────────────────────────────
  ping              Verify Redis connection
  flush             Flush all data (dev only)

  [session]         Auth tokens & DAU bitmaps
  [player]          Profiles & leaderboards
  [comms]           Chat & event log
  [location]        Geo-based matchmaking
  [security]        Rate limiting & threat tracking
  [trading]         Atomic item/currency trading
  [catalog]         JSON catalog & search
`)
}