package main

import (
	"os"
	"testing"
)

func TestConfigParsing(t *testing.T) {
	// Test 1: Single queue mapping
	singleQueueConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  notifications: notifications-queue
  logs: logs-queue`

	os.WriteFile("/tmp/single-queue.yaml", []byte(singleQueueConfig), 0644)
	config1, err := LoadConfig("/tmp/single-queue.yaml")
	if err != nil {
		t.Fatalf("Failed to load single queue config: %v", err)
	}
	
	if len(config1.Mappings["notifications"]) != 1 || config1.Mappings["notifications"][0] != "notifications-queue" {
		t.Errorf("Single queue mapping failed")
	}

	// Test 2: Multiple queue mapping
	multiQueueConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  events:
    - events-queue
    - audit-queue
    - backup-queue`

	os.WriteFile("/tmp/multi-queue.yaml", []byte(multiQueueConfig), 0644)
	config2, err := LoadConfig("/tmp/multi-queue.yaml")
	if err != nil {
		t.Fatalf("Failed to load multi queue config: %v", err)
	}
	
	if len(config2.Mappings["events"]) != 3 {
		t.Errorf("Expected 3 queues, got %d", len(config2.Mappings["events"]))
	}

	// Test 3: Mixed mapping
	mixedConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  notifications: notifications-queue
  events:
    - events-queue
    - audit-queue`

	os.WriteFile("/tmp/mixed.yaml", []byte(mixedConfig), 0644)
	config3, err := LoadConfig("/tmp/mixed.yaml")
	if err != nil {
		t.Fatalf("Failed to load mixed config: %v", err)
	}
	
	if len(config3.Mappings["notifications"]) != 1 {
		t.Errorf("Expected 1 queue for notifications, got %d", len(config3.Mappings["notifications"]))
	}
	if len(config3.Mappings["events"]) != 2 {
		t.Errorf("Expected 2 queues for events, got %d", len(config3.Mappings["events"]))
	}
}

