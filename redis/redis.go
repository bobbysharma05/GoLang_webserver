package redis

import (
    "context"
    "encoding/json"
    "log"

    "github.com/go-redis/redis/v8"
)

type Data struct {
    ID    string `json:"id"`
    Value string `json:"value"`
}

var (
    redisAddr     = "localhost:6379"
    redisPassword = ""
    redisDB       = 0
    redisClient   *redis.Client
)

// InitRedis initializes the Redis client and tests the connection.
func InitRedis() {
    redisClient = redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPassword,
        DB:       redisDB,
    })

    ctx := context.Background()
    _, err := redisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
}

// CloseRedis closes the Redis client.
func CloseRedis() {
    if err := redisClient.Close(); err != nil {
        log.Printf("Error closing Redis client: %v", err)
    }
}

// StoreDataInRedis stores the given data in Redis with the specified key.
func StoreDataInRedis(key string, data Data) {
    ctx := context.Background()

    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Fatalf("Error marshalling data: %v", err)
    }

    err = redisClient.Set(ctx, key, jsonData, 0).Err()
    if err != nil {
        log.Fatalf("Error storing data in Redis: %v", err)
    }
}

// GetDataFromRedis retrieves data from Redis for the given key.
func GetDataFromRedis(key string) (string, error) {
    ctx := context.Background()

    val, err := redisClient.Get(ctx, key).Result()
    if err != nil {
        return "", err
    }

    return val, nil
}