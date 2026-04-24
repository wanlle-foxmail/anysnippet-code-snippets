package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func StreamSSE(writer http.ResponseWriter, request *http.Request, messages <-chan string) error {
	if writer == nil {
		return errors.New("writer is required")
	}
	if request == nil {
		return errors.New("request is required")
	}
	if messages == nil {
		return errors.New("messages channel is required")
	}

	flusher, ok := writer.(http.Flusher)
	if !ok {
		return errors.New("writer does not support streaming")
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	// Flow:
	//   set SSE headers
	//      |
	//      +-> next message -> write "data: ..." block -> flush
	//      `-> channel close or client disconnect -> stop streaming
	for {
		select {
		case <-request.Context().Done():
			return nil
		case message, ok := <-messages:
			if !ok {
				return nil
			}
			if err := writeSSEDataBlock(writer, message); err != nil {
				return fmt.Errorf("write sse event: %w", err)
			}
			flusher.Flush()
		}
	}
}

func writeSSEDataBlock(writer io.Writer, message string) error {
	for _, line := range strings.Split(message, "\n") {
		if _, err := fmt.Fprintf(writer, "data: %s\n", line); err != nil {
			return err
		}
	}
	_, err := io.WriteString(writer, "\n")
	return err
}

func newServer(messages <-chan string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(writer http.ResponseWriter, request *http.Request) {
		if err := StreamSSE(writer, request, messages); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
	return mux
}

func main() {
	messages := make(chan string)
	go func() {
		for index := 1; index <= 3; index++ {
			messages <- fmt.Sprintf("tick-%d", index)
			time.Sleep(250 * time.Millisecond)
		}
		close(messages)
	}()

	http.ListenAndServe(":8080", newServer(messages))
}
