package main

import (
	"log"
	"net/http"

	"webserver/handlers"
	"webserver/kafka"
	"webserver/redis"
)

func main() {
	redis.InitRedis()
	defer redis.CloseRedis()

	log.Println("Starting application...")

	// Start HTTP server
	http.HandleFunc("/produce", handlers.ProduceHandler)
	go func() {
		log.Println("Starting HTTP server on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start Kafka consumer
	log.Println("Starting Kafka consumer...")
	kafka.StartConsumer()
}