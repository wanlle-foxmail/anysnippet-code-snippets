package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func WriteJSONWithETag(writer http.ResponseWriter, request *http.Request, value interface{}) error {
	if writer == nil {
		return errors.New("writer is required")
	}
	if request == nil {
		return errors.New("request is required")
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal response body: %w", err)
	}

	etag := computeETag(payload)
	writer.Header().Set("ETag", etag)
	if etagMatches(request.Header.Get("If-None-Match"), etag) {
		writer.WriteHeader(http.StatusNotModified)
		return nil
	}

	// Flow:
	//   marshal one JSON body
	//      |
	//      +-> If-None-Match matches -> set ETag -> return 304
	//      `-> no match -> set ETag and JSON headers -> write 200 body
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write(append(payload, '\n')); err != nil {
		return fmt.Errorf("write response body: %w", err)
	}
	return nil
}

func computeETag(payload []byte) string {
	hash := sha256.Sum256(payload)
	return `"` + hex.EncodeToString(hash[:]) + `"`
}

func etagMatches(ifNoneMatchHeader string, etag string) bool {
	for _, candidate := range strings.Split(ifNoneMatchHeader, ",") {
		trimmedCandidate := strings.TrimSpace(candidate)
		if trimmedCandidate == "*" || trimmedCandidate == etag {
			return true
		}
	}
	return false
}

func productHandler(writer http.ResponseWriter, request *http.Request) {
	responseBody := map[string]interface{}{
		"product_id": "product-42",
		"status":     "ready",
	}
	if err := WriteJSONWithETag(writer, request, responseBody); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func newServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/products/42", productHandler)
	return mux
}

func main() {
	http.ListenAndServe(":8080", newServer())
}