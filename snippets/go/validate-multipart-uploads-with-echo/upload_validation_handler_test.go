package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestUploadValidationHandlerAcceptsValidTextUpload(t *testing.T) {
	response := performUploadRequest(t, "notes.txt", []byte("hello"), "text/plain")
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	want := map[string]interface{}{
		"filename":     "notes.txt",
		"content_type": "text/plain",
		"size":         float64(5),
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("got %+v; want %+v", payload, want)
	}
}

func TestUploadValidationHandlerRejectsBlankFilename(t *testing.T) {
	response := performUploadRequest(t, "   ", []byte("hello"), "text/plain")
	assertMessage(t, response, http.StatusBadRequest, "filename is required")
}

func TestUploadValidationHandlerRejectsUnsafeFilename(t *testing.T) {
	response := performUploadRequest(t, "..", []byte("hello"), "text/plain")
	assertMessage(t, response, http.StatusBadRequest, "filename is required")
}

func TestUploadValidationHandlerRejectsUnsupportedContentType(t *testing.T) {
	response := performUploadRequest(t, "data.json", []byte("{}"), "application/json")
	assertMessage(t, response, http.StatusUnsupportedMediaType, "unsupported content type")
}

func TestUploadValidationHandlerRejectsSpoofedContentType(t *testing.T) {
	gifBytes := []byte("GIF89a\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\xff\xff\xff,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02L\x01\x00;")
	response := performUploadRequest(t, "image.txt", gifBytes, "text/plain")
	assertMessage(t, response, http.StatusUnsupportedMediaType, "file content does not match declared content type")
}

func TestUploadValidationHandlerRejectsEmptyFiles(t *testing.T) {
	response := performUploadRequest(t, "empty.txt", []byte(""), "text/plain")
	assertMessage(t, response, http.StatusBadRequest, "file is empty")
}

func TestUploadValidationHandlerRejectsFilesLargerThanTheLimit(t *testing.T) {
	response := performUploadRequest(t, "big.txt", bytes.Repeat([]byte("a"), MaxFileSizeBytes+1), "text/plain")
	assertMessage(t, response, http.StatusRequestEntityTooLarge, "file is too large")
}

func TestUploadValidationHandlerAcceptsFilesExactlyAtTheSizeLimit(t *testing.T) {
	response := performUploadRequest(t, "limit.csv", bytes.Repeat([]byte("a"), MaxFileSizeBytes), "text/csv")
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["filename"] != "limit.csv" || payload["content_type"] != "text/csv" || payload["size"] != float64(MaxFileSizeBytes) {
		t.Fatalf("got %+v; want validated size-limit response", payload)
	}
}

func TestUploadValidationHandlerReturnsSanitizedBasename(t *testing.T) {
	response := performUploadRequest(t, "../reports/quarterly.csv", []byte("name,value\n"), "text/csv")
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["filename"] != "quarterly.csv" {
		t.Fatalf("got %+v; want sanitized basename", payload)
	}
}

func TestAllowedContentTypesStaySmallAndExplicit(t *testing.T) {
	want := map[string]struct{}{
		"text/plain": {},
		"text/csv":   {},
	}
	if !reflect.DeepEqual(AllowedContentTypes, want) {
		t.Fatalf("got %+v; want %+v", AllowedContentTypes, want)
	}
}

func performUploadRequest(t *testing.T, filename string, content []byte, contentType string) *httptest.ResponseRecorder {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	partHeaders.Set("Content-Type", contentType)
	partWriter, err := writer.CreatePart(partHeaders)
	if err != nil {
		t.Fatalf("create multipart part: %v", err)
	}
	if _, err := partWriter.Write(content); err != nil {
		t.Fatalf("write multipart content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/upload", body)
	request.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	response := httptest.NewRecorder()
	NewServer().ServeHTTP(response, request)
	return response
}

func assertMessage(t *testing.T, response *httptest.ResponseRecorder, wantStatus int, wantMessage string) {
	t.Helper()
	if response.Code != wantStatus {
		t.Fatalf("expected %d, got %d", wantStatus, response.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["message"] != wantMessage {
		t.Fatalf("got %+v; want message %q", payload, wantMessage)
	}
}
