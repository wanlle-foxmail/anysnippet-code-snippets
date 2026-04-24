package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const jsonlContentType = "application/x-ndjson"

func StreamJSONL(writer http.ResponseWriter, request *http.Request, records <-chan interface{}) error {
	if writer == nil {
		return errors.New("writer is required")
	}
	if request == nil {
		return errors.New("request is required")
	}
	if records == nil {
		return errors.New("records channel is required")
	}

	flusher, ok := writer.(http.Flusher)
	if !ok {
		return errors.New("writer does not support streaming")
	}

	writer.Header().Set("Content-Type", jsonlContentType)
	writer.Header().Set("Cache-Control", "no-cache")

	// Flow:
	//   set JSONL headers
	//      |
	//      +-> next record -> marshal JSON -> write one line -> flush
	//      `-> channel close or client disconnect -> stop streaming
	for {
		select {
		case <-request.Context().Done():
			return nil
		case record, ok := <-records:
			if !ok {
				return nil
			}
			line, err := json.Marshal(record)
			if err != nil {
				return fmt.Errorf("marshal jsonl record: %w", err)
			}
			if _, err := writer.Write(append(line, '\n')); err != nil {
				return fmt.Errorf("write jsonl record: %w", err)
			}
			flusher.Flush()
		}
	}
}

func newServer(records <-chan interface{}) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/records", func(writer http.ResponseWriter, request *http.Request) {
		if err := StreamJSONL(writer, request, records); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
	return mux
}

func main() {
	records := make(chan interface{})
	go func() {
		for index := 1; index <= 3; index++ {
			records <- map[string]interface{}{"id": index, "status": "ready"}
			time.Sleep(250 * time.Millisecond)
		}
		close(records)
	}()

	http.ListenAndServe(":8080", newServer(records))
}