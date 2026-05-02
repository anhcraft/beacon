package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

type Consumer struct {
	client       *pubsub.Client
	subscription *pubsub.Subscription
	handler      *TriggerHandler
	name         string
}

func NewConsumer(ctx context.Context, client *pubsub.Client, name string, config ConsumerConfig) *Consumer {
	sub := client.Subscription(config.PubsubSubscriptionID)
	handler := NewTriggerHandler(name, config)

	return &Consumer{
		client:       client,
		subscription: sub,
		handler:      handler,
		name:         name,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Printf("[%s] Starting consumer on subscription: %s\n", c.name, c.subscription.ID())

	err := c.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("[%s] Received message: %s\n", c.name, string(msg.ID))
		msg.Ack() // Acknowledge immediately
		c.handler.HandleMessage()
	})

	if err != nil {
		return err
	}
	return nil
}
