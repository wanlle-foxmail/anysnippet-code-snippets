package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type CreateBookRequest struct {
	Title string `json:"title"`
}

// Flow:
//
//	wrap the request body with a byte limit
//	   |
//	   +-> one valid JSON value under the limit -> decode and return nil
//	   `-> oversized, malformed, or extra JSON values -> return an error
func DecodeJSONBodyWithLimit(writer http.ResponseWriter, request *http.Request, maxBytes int64, target interface{}) error {
	if writer == nil {
		return errors.New("writer is required")
	}
	if request == nil {
		return errors.New("request is required")
	}
	if target == nil {
		return errors.New("target is required")
	}
	if maxBytes <= 0 {
		return errors.New("max bytes must be greater than 0")
	}

	limitedBody := http.MaxBytesReader(writer, request.Body, maxBytes)
	defer limitedBody.Close()

	decoder := json.NewDecoder(limitedBody)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode JSON body with limit: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain exactly one JSON value")
	}

	return nil
}

func main() {
	request := httptestRequest(`{"title":"Distributed Systems"}`)
	writer := &responseDiscarder{}

	var body CreateBookRequest
	if err := DecodeJSONBodyWithLimit(writer, request, 64, &body); err != nil {
		log.Fatal(err)
	}

	log.Printf("title=%s", body.Title)
}

type responseDiscarder struct{}

func (responseDiscarder) Header() http.Header {
	return make(http.Header)
}

func (responseDiscarder) Write(body []byte) (int, error) {
	return len(body), nil
}

func (responseDiscarder) WriteHeader(statusCode int) {}

func httptestRequest(body string) *http.Request {
	request, err := http.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	return request
}
