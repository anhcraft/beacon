package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func executeCommands(commands []string) error {
	for _, cmdStr := range commands {
		log.Printf("Executing command: %s\n", cmdStr)

		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command execution failed: %s: %w", cmdStr, err)
		}
	}
	return nil
}
