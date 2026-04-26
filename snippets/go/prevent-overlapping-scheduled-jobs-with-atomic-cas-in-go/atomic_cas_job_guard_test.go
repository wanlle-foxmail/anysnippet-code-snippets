package main

import (
	"errors"
	"testing"
)

func TestJobGuardRunsTheJobWhenIdle(t *testing.T) {
	guard := &JobGuard{}
	ran, err := guard.Run(func() error { return nil })
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !ran {
		t.Fatal("expected the job to run")
	}
}

func TestJobGuardSkipsOverlappingRuns(t *testing.T) {
	guard := &JobGuard{}
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)

	go func() {
		_, err := guard.Run(func() error {
			close(started)
			<-release
			return nil
		})
		done <- err
	}()

	<-started
	ran, err := guard.Run(func() error { return nil })
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if ran {
		t.Fatal("expected the overlapping run to be skipped")
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestJobGuardClearsTheRunningFlagAfterSuccess(t *testing.T) {
	guard := &JobGuard{}
	if _, err := guard.Run(func() error { return nil }); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if guard.Running() {
		t.Fatal("expected running to be false after success")
	}
}

func TestJobGuardClearsTheRunningFlagAfterErrors(t *testing.T) {
	guard := &JobGuard{}
	wantErr := errors.New("job failed")
	ran, err := guard.Run(func() error { return wantErr })
	if !ran {
		t.Fatal("expected the job to run")
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
	if guard.Running() {
		t.Fatal("expected running to be false after an error")
	}
}

func TestJobGuardClearsTheRunningFlagAfterPanics(t *testing.T) {
	guard := &JobGuard{}
	defer func() {
		recovered := recover()
		if recovered != "boom" {
			t.Fatalf("expected panic boom, got %v", recovered)
		}
		if guard.Running() {
			t.Fatal("expected running to be false after a panic")
		}
	}()

	_, _ = guard.Run(func() error {
		panic("boom")
	})
}

func TestJobGuardRejectsInvalidInput(t *testing.T) {
	var nilGuard *JobGuard
	if _, err := nilGuard.Run(func() error { return nil }); err == nil {
		t.Fatal("expected error for a nil guard")
	}

	guard := &JobGuard{}
	if _, err := guard.Run(nil); err == nil {
		t.Fatal("expected error for a nil job")
	}
}
