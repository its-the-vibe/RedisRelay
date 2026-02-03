# RedisRelay

A simple Go service that subscribes to Redis channels and publishes messages to Redis lists.

## Overview

RedisRelay is a lightweight service that acts as a bridge between Redis Pub/Sub channels and Redis lists. It subscribes to one or more Redis channels and automatically pushes received messages to corresponding Redis lists (queues), enabling asynchronous message processing patterns.

## Features

- Written in Go 1.24
- Minimal Docker image using scratch base (multi-stage build)
- Configurable channel-to-queue mappings via YAML
- Connects to external Redis servers
- Graceful shutdown handling
- Simple and efficient message relay

## Configuration

Create a `config.yaml` file based on the provided `config.yaml.example`:

```yaml
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

mappings:
  # Single queue mapping
  notifications: notifications-queue
  logs: logs-queue
  
  # Multiple queue mapping (fan-out pattern)
  events:
    - events-queue
    - audit-queue
```

### Configuration Options

- `redis.host`: Redis server hostname
- `redis.port`: Redis server port
- `redis.password`: Redis authentication password (optional)
- `redis.db`: Redis database number
- `mappings`: Map of channel names (keys) to queue names (values)
  - **Single queue**: Use a string value for one-to-one mapping
  - **Multiple queues**: Use an array of strings for one-to-many mapping (fan-out pattern)

## Running Locally

### Prerequisites

- Go 1.24 or later
- Access to a Redis server

### Build and Run

1. Copy the example configuration:
   ```bash
   cp config.yaml.example config.yaml
   ```

2. Edit `config.yaml` with your Redis connection details

3. Build the application:
   ```bash
   go build -o redisrelay
   ```

4. Run the application:
   ```bash
   ./redisrelay
   ```

   Or specify a custom config path:
   ```bash
   CONFIG_PATH=/path/to/config.yaml ./redisrelay
   ```

## Running with Docker

### Build the Docker image

```bash
docker build -t redisrelay:latest .
```

### Run with Docker

```bash
docker run -v $(pwd)/config.yaml:/config.yaml:ro \
  -e CONFIG_PATH=/config.yaml \
  redisrelay:latest
```

### Run with Docker Compose

1. Copy the example configuration:
   ```bash
   cp config.yaml.example config.yaml
   ```

2. Edit `config.yaml` with your Redis connection details

3. Start the service:
   ```bash
   docker-compose up -d
   ```

4. View logs:
   ```bash
   docker-compose logs -f
   ```

5. Stop the service:
   ```bash
   docker-compose down
   ```

## How It Works

1. The service reads the configuration file to get Redis connection details and channel mappings
2. It connects to the Redis server and subscribes to all configured channels
3. When a message is published to a subscribed channel, the service receives it
4. The message payload is pushed to the corresponding Redis list(s) using RPUSH
   - For single queue mappings, the message is pushed to one queue
   - For multiple queue mappings (fan-out), the message is pushed to all configured queues
5. The process continues until the service is terminated

## Testing

You can test the service using Redis CLI:

```bash
# Publish a message to a channel
redis-cli PUBLISH events "test message"

# Check if the message was added to the queue
redis-cli LRANGE events-queue 0 -1
```

## License

MIT
