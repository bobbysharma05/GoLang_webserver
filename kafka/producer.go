package kafka

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"webserver/redis"
)

// ProduceMessage sends a message to the Kafka topic.
func ProduceMessage(data redis.Data) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
	})
	if err != nil {
		return err
	}
	defer producer.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: int32(kafka.PartitionAny)},
		Value:          jsonData,
		Key:            []byte(data.ID),
	}, nil)
	if err != nil {
		return err
	}

	producer.Flush(15 * 1000)
	log.Printf("Produced message to Kafka: %+v", data)
	return nil
}