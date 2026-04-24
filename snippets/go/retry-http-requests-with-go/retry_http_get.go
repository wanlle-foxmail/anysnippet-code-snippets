package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

var sleep = time.Sleep

func RetryHTTPGet(client *http.Client, url string, maxAttempts int, delay time.Duration) (*http.Response, error) {
	// Flow: send GET -> success or non-retryable response returns immediately -> retryable status or transport error sleeps and retries -> last attempt returns final response or error
	if client == nil {
		return nil, errors.New("client is required")
	}
	if url == "" {
		return nil, errors.New("url is required")
	}
	if maxAttempts <= 0 {
		return nil, errors.New("max attempts must be greater than 0")
	}
	if delay < 0 {
		return nil, errors.New("delay must be greater than or equal to 0")
	}

	var lastTransportError error

	for attemptNumber := 1; attemptNumber <= maxAttempts; attemptNumber++ {
		response, err := client.Get(url)
		if err != nil {
			lastTransportError = err
			if attemptNumber == maxAttempts {
				return nil, fmt.Errorf("get %s: %w", url, err)
			}
			if delay > 0 {
				sleep(delay)
			}
			continue
		}

		if !shouldRetryStatus(response.StatusCode) || attemptNumber == maxAttempts {
			return response, nil
		}

		if response.Body != nil {
			_ = response.Body.Close()
		}
		if delay > 0 {
			sleep(delay)
		}
	}

	if lastTransportError != nil {
		return nil, fmt.Errorf("get %s: %w", url, lastTransportError)
	}

	return nil, errors.New("retry loop ended unexpectedly")
}

func shouldRetryStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

func main() {
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := RetryHTTPGet(client, "https://example.com/health", 3, 500*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	log.Printf("status=%d", response.StatusCode)
}
