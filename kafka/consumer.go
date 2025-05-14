package kafka

import (
    "encoding/json"
    "log"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "webserver/redis"
)

var (
    kafkaGroup  = "redis-consumer-group"
)

// StartConsumer starts a Kafka consumer that reads messages and stores them in Redis.
func StartConsumer() {
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": kafkaBroker,
        "group.id":          kafkaGroup,
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        log.Fatalf("Failed to create Kafka consumer: %v", err)
    }
    defer consumer.Close()

    err = consumer.SubscribeTopics([]string{kafkaTopic}, nil)
    if err != nil {
        log.Fatalf("Failed to subscribe to topic: %v", err)
    }

    log.Println("Starting Kafka consumer...")   
    for {
        msg, err := consumer.ReadMessage(-1)
        if err != nil {
            log.Printf("Consumer error: %v", err)
            continue
        }

        var data redis.Data
        if err := json.Unmarshal(msg.Value, &data); err != nil {
            log.Printf("Error unmarshalling Kafka message: %v", err)
            continue
        }

        redisKey := string(msg.Key)
        if redisKey == "" {
            redisKey = data.ID
        }

        redis.StoreDataInRedis(redisKey, data)
        log.Printf("Stored data in Redis with key %s: %+v", redisKey, data)

        retrievedData, err := redis.GetDataFromRedis(redisKey)
        if err != nil {
            log.Printf("Error retrieving data from Redis: %v", err)
            continue
        }
        log.Printf("Retrieved data from Redis: %s", retrievedData)
    }
}