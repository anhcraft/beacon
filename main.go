package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"cloud.google.com/go/pubsub"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to the configuration file")
	flag.Parse()

	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := pubsub.NewClient(ctx, config.GCPProjectID)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer client.Close()

	var wg sync.WaitGroup

	for name, consumerConfig := range config.Consumers {
		consumer := NewConsumer(ctx, client, name, consumerConfig)
		wg.Add(1)
		go func(c *Consumer, n string) {
			defer wg.Done()
			if err := c.Start(ctx); err != nil {
				log.Printf("[%s] Consumer stopped with error: %v\n", n, err)
			} else {
				log.Printf("[%s] Consumer stopped cleanly.\n", n)
			}
		}(consumer, name)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Beacon started. Waiting for messages...")

	sig := <-sigChan
	log.Printf("Received signal %v, initiating shutdown...\n", sig)

	// Cancel the context to stop all consumers
	cancel()

	// Wait for consumers to finish processing current messages
	wg.Wait()
	log.Println("Beacon shutdown complete.")
}
