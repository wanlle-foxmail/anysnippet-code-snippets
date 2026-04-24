package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func DecodeStrictJSON(reader io.Reader, target interface{}) error {
	// Flow:
	//   decode one JSON value with unknown-field rejection enabled
	//      |
	//      +-> one valid object -> fill target and return nil
	//      `-> unknown fields, malformed JSON, or extra values -> return an error
	if reader == nil {
		return errors.New("reader is required")
	}
	if target == nil {
		return errors.New("target is required")
	}

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode strict JSON: %w", err)
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain exactly one JSON value")
	}

	return nil
}

func main() {
	requestBody := strings.NewReader(`{"display_name":"Ada","email":"ada@example.com"}`)
	var request CreateUserRequest
	if err := DecodeStrictJSON(requestBody, &request); err != nil {
		log.Fatal(err)
	}

	log.Printf("request=%+v", request)
}