package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Event struct {
	ID   int               `json:"id"`
	Tags []string          `json:"tags"`
	Meta map[string]string `json:"meta"`
}

// Flow:
//
//	open one top-level JSON array token
//	   |
//	   +-> array items remain -> decode and handle one item at a time
//	   `-> malformed or non-array input -> return an error without loading the full array at once
func ReadJSONArrayItems[T any](reader io.Reader, handle func(T) error) error {
	if reader == nil {
		return errors.New("reader is required")
	}
	if handle == nil {
		return errors.New("handler is required")
	}

	decoder := json.NewDecoder(reader)
	startToken, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("read array start: %w", err)
	}
	delim, ok := startToken.(json.Delim)
	if !ok || delim != '[' {
		return errors.New("input must start with a JSON array")
	}

	for decoder.More() {
		var item T
		if err := decoder.Decode(&item); err != nil {
			return fmt.Errorf("decode array item: %w", err)
		}
		if err := handle(item); err != nil {
			return fmt.Errorf("handle array item: %w", err)
		}
	}

	endToken, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("read array end: %w", err)
	}
	delim, ok = endToken.(json.Delim)
	if !ok || delim != ']' {
		return errors.New("input must end with a JSON array")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("input must contain exactly one JSON array")
	}

	return nil
}

func main() {
	reader := strings.NewReader(`[{"id":1,"tags":["blue"],"meta":{"source":"api"}},{"id":2,"tags":["green"],"meta":{"source":"cache"}}]`)
	if err := ReadJSONArrayItems(reader, func(item Event) error {
		log.Printf("event=%+v", item)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
