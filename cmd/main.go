package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sayandip18/redis-game-engine/internal/client"
	"github.com/sayandip18/redis-game-engine/internal/session"
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
		case "session":
			if len(os.Args) < 3 {
				fmt.Println("Usage: session <create|validate|dau> [args...]")
				return
			}
			switch os.Args[2] {
			case "create":
				if len(os.Args) < 5 {
					fmt.Println("Usage: session create <userID> <ttl>")
					return
				}
				token, err := session.CreateSession(os.Args[3], os.Args[4])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("token:", token)
			case "validate":
				if len(os.Args) < 4 {
					fmt.Println("Usage: session validate <token>")
					return
				}
				userID, err := session.ValidateSession(os.Args[3])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("userID:", userID)
			case "dau":
				if len(os.Args) < 4 || os.Args[3] != "today" {
					fmt.Println("Usage: session dau today")
					return
				}
				today := time.Now().Format("2006-01-02")
				count, err := session.GetDAU(today)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("DAU %s: %d\n", today, count)
			default:
				fmt.Printf("Unknown session subcommand: %s\n", os.Args[2])
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