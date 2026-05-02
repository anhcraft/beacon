package main

import (
	"log"
	"sync"
	"time"
)

type TriggerHandler struct {
	name            string
	config          ConsumerConfig
	timer           *time.Timer
	lastMessageTime time.Time
	mu              sync.Mutex
}

func NewTriggerHandler(name string, config ConsumerConfig) *TriggerHandler {
	return &TriggerHandler{
		name:   name,
		config: config,
	}
}

func (h *TriggerHandler) HandleMessage() {
	if !h.config.Deduplication.Enabled {
		log.Printf("[%s] Deduplication disabled, executing immediately...\n", h.name)
		h.execute()
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.lastMessageTime = time.Now()

	if h.timer != nil {
		log.Printf("[%s] Message received within deduplication window. Timer will be extended.\n", h.name)
	} else {
		log.Printf("[%s] Message received. Starting deduplication window of %v...\n", h.name, h.config.Deduplication.TimeWindow)
		h.timer = time.AfterFunc(h.config.Deduplication.TimeWindow, h.checkTimer)
	}
}

func (h *TriggerHandler) checkTimer() {
	h.mu.Lock()
	elapsed := time.Since(h.lastMessageTime)
	remaining := h.config.Deduplication.TimeWindow - elapsed

	if remaining > 0 {
		// Not enough time has passed, meaning we received another message. Reschedule.
		h.timer = time.AfterFunc(remaining, h.checkTimer)
		h.mu.Unlock()
		return
	}

	// Time elapsed without new messages.
	h.timer = nil
	h.mu.Unlock()

	log.Printf("[%s] Deduplication window elapsed. Executing commands...\n", h.name)
	h.execute()
}

func (h *TriggerHandler) execute() {
	if err := executeCommands(h.config.TriggerCommands); err != nil {
		log.Printf("[%s] Error executing commands: %v\n", h.name, err)
	}
}
