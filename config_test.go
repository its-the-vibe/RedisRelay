package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigParsing(t *testing.T) {
	dir := t.TempDir()

	// Test 1: Single queue mapping
	singleQueueConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  notifications: notifications-queue
  logs: logs-queue`

	singleQueueFile := filepath.Join(dir, "single-queue.yaml")
	os.WriteFile(singleQueueFile, []byte(singleQueueConfig), 0644)
	config1, err := LoadConfig(singleQueueFile)
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

	multiQueueFile := filepath.Join(dir, "multi-queue.yaml")
	os.WriteFile(multiQueueFile, []byte(multiQueueConfig), 0644)
	config2, err := LoadConfig(multiQueueFile)
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

	mixedFile := filepath.Join(dir, "mixed.yaml")
	os.WriteFile(mixedFile, []byte(mixedConfig), 0644)
	config3, err := LoadConfig(mixedFile)
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

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Error("Expected error for missing config file, got nil")
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	invalidYAML := `redis: {invalid yaml`
	path := filepath.Join(dir, "invalid.yaml")
	os.WriteFile(path, []byte(invalidYAML), 0644)

	_, err := LoadConfig(path)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoadConfigRedisPasswordFromEnv(t *testing.T) {
	dir := t.TempDir()
	configWithoutPassword := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  events: events-queue`

	path := filepath.Join(dir, "no-password.yaml")
	os.WriteFile(path, []byte(configWithoutPassword), 0644)
	t.Setenv("REDIS_PASSWORD", "secret123")

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if config.Redis.Password != "secret123" {
		t.Errorf("Expected password 'secret123', got '%s'", config.Redis.Password)
	}
}

func TestLoadConfigPasswordInFileNotOverridden(t *testing.T) {
	dir := t.TempDir()
	configWithPassword := `redis:
  host: localhost
  port: 6379
  password: filepassword
  db: 0

mappings:
  events: events-queue`

	path := filepath.Join(dir, "with-password.yaml")
	os.WriteFile(path, []byte(configWithPassword), 0644)
	t.Setenv("REDIS_PASSWORD", "envpassword")

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if config.Redis.Password != "filepassword" {
		t.Errorf("Expected password 'filepassword' (from file), got '%s'", config.Redis.Password)
	}
}

func TestNormalizeMappingsInvalidType(t *testing.T) {
	dir := t.TempDir()
	// A mapping value that is neither a string nor a string array should error
	invalidTypeConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  events: 12345`

	path := filepath.Join(dir, "invalid-type.yaml")
	os.WriteFile(path, []byte(invalidTypeConfig), 0644)
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("Expected error for invalid mapping type, got nil")
	}
}

func TestNormalizeMappingsEmptyArrayError(t *testing.T) {
	dir := t.TempDir()
	emptyArrayConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  events: []`

	path := filepath.Join(dir, "empty-array.yaml")
	os.WriteFile(path, []byte(emptyArrayConfig), 0644)
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("Expected error for empty queue array, got nil")
	}
}

