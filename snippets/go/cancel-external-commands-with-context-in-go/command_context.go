package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Flow: start the command with context -> capture combined output -> prefer ctx.Err() when cancellation stops the process.
func RunCommandWithContext(ctx context.Context, name string, args ...string) ([]byte, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("command name is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	command := exec.CommandContext(ctx, name, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return output, fmt.Errorf("run %s: %w", name, ctxErr)
		}
		return output, fmt.Errorf("run %s: %w", name, err)
	}

	return output, nil
}

func main() {
	output, err := RunCommandWithContext(context.Background(), "go", "version")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("output=%s", strings.TrimSpace(string(output)))
}
