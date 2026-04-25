package main

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestSignalContextRejectsEmptySignalList(t *testing.T) {
	if _, _, err := SignalContext(context.Background()); err == nil {
		t.Fatal("expected error when no signals are provided")
	}
}

func TestSignalContextUsesBackgroundForNilParent(t *testing.T) {
	ctx, stop, err := SignalContext(nil, os.Interrupt)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer stop()
	if ctx == nil {
		t.Fatal("expected a derived context")
	}
}

func TestSignalContextHonorsParentCancellation(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	ctx, stop, err := SignalContext(parent, os.Interrupt)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer stop()

	cancelParent()

	select {
	case <-ctx.Done():
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for parent cancellation to reach the derived context")
	}
}

func TestSignalContextStopCancelsContext(t *testing.T) {
	ctx, stop, err := SignalContext(context.Background(), os.Interrupt)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	stop()

	select {
	case <-ctx.Done():
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for stop to cancel the derived context")
	}
}

func TestSignalContextCancelsOnInterruptSignal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Process.Signal(os.Interrupt) is not reliably supported on Windows test processes")
	}

	command := exec.Command(os.Args[0], "-test.run=TestSignalContextHelperProcess")
	command.Env = append(os.Environ(), "GO_WANT_SIGNAL_HELPER=1")
	if err := command.Start(); err != nil {
		t.Fatalf("start helper process: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	if err := command.Process.Signal(os.Interrupt); err != nil {
		_ = command.Process.Kill()
		t.Fatalf("send interrupt: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- command.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("helper process exited with error: %v", err)
		}
	case <-time.After(2 * time.Second):
		_ = command.Process.Kill()
		t.Fatal("timed out waiting for helper process to exit after interrupt")
	}
}

func TestSignalContextHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_SIGNAL_HELPER") != "1" {
		return
	}

	ctx, stop, err := SignalContext(context.Background(), os.Interrupt)
	if err != nil {
		os.Exit(2)
	}
	defer stop()

	<-ctx.Done()
	os.Exit(0)
}
