package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
)

func SignalContext(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc, error) {
	if len(signals) == 0 {
		return nil, nil, errors.New("at least one signal is required")
	}
	if parent == nil {
		parent = context.Background()
	}

	ctx, stop := signal.NotifyContext(parent, signals...)
	return ctx, stop, nil
}

func main() {
	ctx, stop, err := SignalContext(context.Background(), os.Interrupt)
	if err != nil {
		log.Fatal(err)
	}
	defer stop()

	log.Println("waiting for interrupt signal")
	<-ctx.Done()
	log.Println("signal received, shutting down")
}
