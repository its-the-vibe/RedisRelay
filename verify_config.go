//go:build ignore
// +build ignore

package main

import (
"fmt"
"os"
)

func main() {
fmt.Println("Verifying configuration parsing...")

// Test backward compatibility
singleConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  notifications: notifications-queue
  logs: logs-queue`

os.WriteFile("/tmp/single-test.yaml", []byte(singleConfig), 0644)
config1, err := LoadConfig("/tmp/single-test.yaml")
if err != nil {
fmt.Printf("❌ ERROR: %v\n", err)
os.Exit(1)
}

fmt.Println("\n✅ Test 1: Backward Compatibility (Single Queue)")
for channel, queues := range config1.Mappings {
fmt.Printf("  Channel '%s' -> Queues: %v\n", channel, queues)
}

// Test new functionality
multiConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  events:
    - events-queue
    - audit-queue
    - backup-queue`

os.WriteFile("/tmp/multi-test.yaml", []byte(multiConfig), 0644)
config2, err := LoadConfig("/tmp/multi-test.yaml")
if err != nil {
fmt.Printf("❌ ERROR: %v\n", err)
os.Exit(1)
}

fmt.Println("\n✅ Test 2: One-to-Many Mapping")
for channel, queues := range config2.Mappings {
fmt.Printf("  Channel '%s' -> Queues: %v\n", channel, queues)
}

// Test mixed configuration
mixedConfig := `redis:
  host: localhost
  port: 6379
  db: 0

mappings:
  notifications: notifications-queue
  logs: logs-queue
  events:
    - events-queue
    - audit-queue`

os.WriteFile("/tmp/mixed-test.yaml", []byte(mixedConfig), 0644)
config3, err := LoadConfig("/tmp/mixed-test.yaml")
if err != nil {
fmt.Printf("❌ ERROR: %v\n", err)
os.Exit(1)
}

fmt.Println("\n✅ Test 3: Mixed Configuration")
for channel, queues := range config3.Mappings {
fmt.Printf("  Channel '%s' -> Queues: %v\n", channel, queues)
}

fmt.Println("\n🎉 All configuration tests passed successfully!")
}
