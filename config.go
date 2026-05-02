package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GCPProjectID string                    `yaml:"gcp-project-id"`
	Consumers    map[string]ConsumerConfig `yaml:"consumers"`
}

type ConsumerConfig struct {
	PubsubSubscriptionID string              `yaml:"pubsub-subscription-id"`
	Deduplication        DeduplicationConfig `yaml:"deduplication"`
	TriggerCommands      []string            `yaml:"trigger-commands"`
}

type DeduplicationConfig struct {
	Enabled    bool          `yaml:"enabled"`
	TimeWindow time.Duration `yaml:"time-window"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}
