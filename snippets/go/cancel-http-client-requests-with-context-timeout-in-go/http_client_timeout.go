package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Flow: derive a request-scoped timeout -> build the GET request with context -> return the response or the cancellation error.
func GetWithContextTimeout(ctx context.Context, client *http.Client, url string, timeout time.Duration) (*http.Response, error) {
	if client == nil {
		return nil, errors.New("client is required")
	}
	if url == "" {
		return nil, errors.New("url is required")
	}
	if timeout <= 0 {
		return nil, errors.New("timeout must be greater than 0")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	requestContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(requestContext, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %s: %w", url, err)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", url, err)
	}

	return response, nil
}

func main() {
	client := &http.Client{}
	response, err := GetWithContextTimeout(context.Background(), client, "https://example.com/health", 2*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	log.Printf("status=%d", response.StatusCode)
}
