package main

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRunCommandWithContextReturnsCombinedOutput(t *testing.T) {
	t.Setenv("GO_WANT_RUN_COMMAND_HELPER", "1")
	t.Setenv("RUN_COMMAND_HELPER_MODE", "success")

	output, err := RunCommandWithContext(context.Background(), os.Args[0], "-test.run=TestRunCommandHelperProcess")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if strings.TrimSpace(string(output)) != "ok" {
		t.Fatalf("expected output ok, got %q", string(output))
	}
}

func TestRunCommandWithContextReturnsCombinedErrorOutput(t *testing.T) {
	t.Setenv("GO_WANT_RUN_COMMAND_HELPER", "1")
	t.Setenv("RUN_COMMAND_HELPER_MODE", "error")

	output, err := RunCommandWithContext(context.Background(), os.Args[0], "-test.run=TestRunCommandHelperProcess")
	if err == nil {
		t.Fatal("expected a command error")
	}
	if strings.TrimSpace(string(output)) != "boom" {
		t.Fatalf("expected stderr output boom, got %q", string(output))
	}
}

func TestRunCommandWithContextCancelsLongRunningCommand(t *testing.T) {
	t.Setenv("GO_WANT_RUN_COMMAND_HELPER", "1")
	t.Setenv("RUN_COMMAND_HELPER_MODE", "block")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := RunCommandWithContext(ctx, os.Args[0], "-test.run=TestRunCommandHelperProcess")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
}

func TestRunCommandWithContextHonorsParentCancellation(t *testing.T) {
	t.Setenv("GO_WANT_RUN_COMMAND_HELPER", "1")
	t.Setenv("RUN_COMMAND_HELPER_MODE", "success")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := RunCommandWithContext(ctx, os.Args[0], "-test.run=TestRunCommandHelperProcess")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestRunCommandWithContextRejectsEmptyCommandName(t *testing.T) {
	if _, err := RunCommandWithContext(context.Background(), "   "); err == nil {
		t.Fatal("expected error for empty command name")
	}
}

func TestRunCommandWithContextUsesBackgroundForNilContext(t *testing.T) {
	t.Setenv("GO_WANT_RUN_COMMAND_HELPER", "1")
	t.Setenv("RUN_COMMAND_HELPER_MODE", "success")

	output, err := RunCommandWithContext(nil, os.Args[0], "-test.run=TestRunCommandHelperProcess")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if strings.TrimSpace(string(output)) != "ok" {
		t.Fatalf("expected output ok, got %q", string(output))
	}
}

func TestRunCommandHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_RUN_COMMAND_HELPER") != "1" {
		return
	}

	switch os.Getenv("RUN_COMMAND_HELPER_MODE") {
	case "success":
		_, _ = os.Stdout.WriteString("ok")
		os.Exit(0)
	case "error":
		_, _ = os.Stderr.WriteString("boom")
		os.Exit(2)
	case "block":
		time.Sleep(10 * time.Second)
		os.Exit(0)
	default:
		os.Exit(3)
	}
}
