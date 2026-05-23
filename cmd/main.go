package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sayandip18/redis-game-engine/internal/client"
	"github.com/sayandip18/redis-game-engine/internal/player"
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
		case "player":
			if len(os.Args) < 3 {
				fmt.Println("Usage: player <create|profile|top10|rank> [args...]")
				return
			}
			switch os.Args[2] {
			case "create":
				if len(os.Args) < 5 {
					fmt.Println("Usage: player create <userID> <username>")
					return
				}
				p, err := player.CreatePlayer(os.Args[3], os.Args[4])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("userID:   %s\nusername: %s\nMMR:      %d\ncreated:  %s\n",
					p.UserID, p.Username, p.MMR, p.CreatedAt.Format(time.RFC3339))
			case "profile":
				if len(os.Args) < 4 {
					fmt.Println("Usage: player profile <userID>")
					return
				}
				p, err := player.GetProfile(os.Args[3])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("userID:   %s\nusername: %s\nMMR:      %d\ncreated:  %s\n",
					p.UserID, p.Username, p.MMR, p.CreatedAt.Format(time.RFC3339))
				if len(p.Stats) > 0 {
					fmt.Println("stats:")
					for k, v := range p.Stats {
						fmt.Printf("  %s: %d\n", k, v)
					}
				}
			case "top10":
				entries, err := player.GetGlobalTop(10)
				if err != nil {
					log.Fatal(err)
				}
				for _, e := range entries {
					fmt.Printf("#%-3d %-20s %d\n", e.Rank, e.PlayerID, e.MMR)
				}
			case "rank":
				if len(os.Args) < 4 {
					fmt.Println("Usage: player rank <userID>")
					return
				}
				rank, err := player.GetRank(os.Args[3])
				if err != nil {
					log.Fatal(err)
				}
				p, err := player.GetProfile(os.Args[3])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("rank: #%d  MMR: %d\n", rank, p.MMR)
			default:
				fmt.Printf("Unknown player subcommand: %s\n", os.Args[2])
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