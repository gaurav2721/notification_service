package consumers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gaurav2721/notification-service/external_services/kafka"
)

// ExampleUsage demonstrates how to use the consumer manager in your main application
func ExampleUsage() {
	// Initialize your Kafka service (this would be your actual implementation)
	var kafkaService kafka.KafkaService
	// kafkaService = your_kafka_implementation.NewKafkaService()

	// Create consumer manager from environment configuration
	consumerManager := NewConsumerManagerFromEnv(kafkaService)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the consumer manager
	if err := consumerManager.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize consumer manager: %v", err)
	}

	// Start all worker pools
	if err := consumerManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start consumer manager: %v", err)
	}

	log.Println("Consumer manager started successfully")

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, stopping consumer manager...")

	// Stop the consumer manager gracefully
	if err := consumerManager.Stop(); err != nil {
		log.Printf("Error stopping consumer manager: %v", err)
	}

	log.Println("Consumer manager stopped successfully")
}

// ExampleWithCustomConfiguration demonstrates how to use the consumer manager with custom configuration
func ExampleWithCustomConfiguration() {
	// Initialize your Kafka service
	var kafkaService kafka.KafkaService
	// kafkaService = your_kafka_implementation.NewKafkaService()

	// Create custom configuration
	config := ConsumerConfig{
		EmailWorkerCount:       10, // More workers for email
		SlackWorkerCount:       2,  // Fewer workers for slack
		IOSPushWorkerCount:     5,  // Medium workers for iOS
		AndroidPushWorkerCount: 5,  // Medium workers for Android
		KafkaService:           kafkaService,
	}

	// Create consumer manager with custom configuration
	consumerManager := NewConsumerManager(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize and start
	if err := consumerManager.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize consumer manager: %v", err)
	}

	if err := consumerManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start consumer manager: %v", err)
	}

	// Example: Get status of all worker pools
	status := consumerManager.GetStatus()
	for notificationType, isRunning := range status {
		log.Printf("Worker pool for %s: %v", notificationType, isRunning)
	}

	// Example: Get a specific worker pool
	if emailPool, err := consumerManager.GetWorkerPool(EmailNotification); err == nil {
		log.Printf("Email worker pool has %d workers", emailPool.GetWorkerCount())
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err := consumerManager.Stop(); err != nil {
		log.Printf("Error stopping consumer manager: %v", err)
	}
}
