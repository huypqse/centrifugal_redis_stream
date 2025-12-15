package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Load Lua script
	luaScript, err := os.ReadFile("lua/update_leaderboard.lua")
	if err != nil {
		log.Fatalf("Failed to read Lua script: %v", err)
	}

	// Register Lua script
	updateLeaderboard := redis.NewScript(string(luaScript))

	// Leader names
	leaderNames := []string{
		"Alice", "Bob", "Charlie", "David", "Eve",
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Main loop
	for {
		leader := leaderNames[rand.Intn(len(leaderNames))]
		increment := rand.Intn(10) + 1
		channel := "leaderboard"

		// Execute Lua script
		_, err := updateLeaderboard.Run(
			ctx,
			rdb,
			[]string{"leaderboard", "leaderboard-state", "leaderboard-stream"},
			leader,
			increment,
			channel,
		).Result()

		if err != nil {
			log.Printf("Error executing Lua script: %v", err)
		} else {
			fmt.Printf("Updated %s with +%d points\n", leader, increment)
		}

		time.Sleep(200 * time.Millisecond)
	}
}
