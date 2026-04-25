package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

// Flow:
//
//	iterate items
//	   |
//	   +-> context canceled -> return ctx.Err() before starting the next item
//	   +-> process returns error -> wrap the item and stop
//	   `-> success -> continue until the batch is done
func ProcessBatchWithContext(ctx context.Context, items []string, process func(context.Context, string) error) error {
	if process == nil {
		return errors.New("process function is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	for _, item := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := process(ctx, item); err != nil {
			return fmt.Errorf("process %q: %w", item, err)
		}
	}

	return nil
}

func main() {
	items := []string{"invoice-1", "invoice-2", "invoice-3"}
	err := ProcessBatchWithContext(context.Background(), items, func(ctx context.Context, item string) error {
		log.Printf("processing %s", item)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
