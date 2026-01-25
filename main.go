package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	defer redisClient.Close()

	// Test connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	// Create a pubsub instance
	pubsub := redisClient.PSubscribe(ctx)
	defer pubsub.Close()

	// Subscribe to all configured channels
	channels := make([]string, 0, len(config.Mappings))
	for channel := range config.Mappings {
		channels = append(channels, channel)
	}

	if len(channels) == 0 {
		log.Fatal("No channel mappings configured")
	}

	if err := pubsub.Subscribe(ctx, channels...); err != nil {
		log.Fatalf("Failed to subscribe to channels: %v", err)
	}

	log.Printf("Subscribed to channels: %v", channels)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start message processing
	go processMessages(ctx, pubsub, redisClient, config.Mappings)

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down gracefully...")
}

// processMessages processes incoming messages from subscribed channels
func processMessages(ctx context.Context, pubsub *redis.PubSub, client *redis.Client, mappings map[string]string) {
	ch := pubsub.Channel()

	for msg := range ch {
		targetQueue, ok := mappings[msg.Channel]
		if !ok {
			log.Printf("No mapping found for channel: %s", msg.Channel)
			continue
		}

		// Push the message payload to the target Redis list
		if err := client.RPush(ctx, targetQueue, msg.Payload).Err(); err != nil {
			log.Printf("Failed to push message to queue %s: %v", targetQueue, err)
		} else {
			log.Printf("Message from channel %s pushed to queue %s", msg.Channel, targetQueue)
		}
	}
}
