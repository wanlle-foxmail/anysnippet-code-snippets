package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postBookRequest(t *testing.T, payload string) (*httptest.ResponseRecorder, map[string]interface{}) {
	t.Helper()

	e := NewServer()
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	return rec, body
}

func validationErrors(t *testing.T, body map[string]interface{}) map[string]interface{} {
	t.Helper()

	errorsMap, ok := body["errors"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected errors map, got %T", body["errors"])
	}

	return errorsMap
}

func TestCreateBookHandlerRunsNextStepForValidInput(t *testing.T) {
	originalProcessBook := ProcessBook
	defer func() { ProcessBook = originalProcessBook }()

	called := false
	var received Book
	ProcessBook = func(book Book) map[string]interface{} {
		called = true
		received = book
		return map[string]interface{}{
			"message": "Book created successfully.",
		}
	}

	rec, body := postBookRequest(t, `{"name":"Clean Code","page":464,"type":"reference","description":"A practical guide to writing better code."}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	if !called {
		t.Fatal("expected ProcessBook to be called")
	}
	if received.Name != "Clean Code" {
		t.Fatalf("expected name Clean Code, got %q", received.Name)
	}
	if body["message"] != "Book created successfully." {
		t.Fatalf("expected success message, got %v", body["message"])
	}
}

func TestCreateBookHandlerReturnsFriendlyMessageForMissingName(t *testing.T) {
	rec, body := postBookRequest(t, `{"name":"   ","page":120,"type":"fiction","description":"A friendly story for weekend reading."}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if body["message"] != "Please correct the highlighted fields and try again." {
		t.Fatalf("expected friendly validation message, got %v", body["message"])
	}
	if validationErrors(t, body)["name"] != "Please enter a book name." {
		t.Fatalf("expected friendly name message, got %v", validationErrors(t, body)["name"])
	}
}

func TestCreateBookHandlerReturnsFriendlyMessageForInvalidPage(t *testing.T) {
	rec, body := postBookRequest(t, `{"name":"Go Basics","page":0,"type":"reference","description":"A short guide for learning Go basics."}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if validationErrors(t, body)["page"] != "Please enter a page count greater than 0." {
		t.Fatalf("expected friendly page message, got %v", validationErrors(t, body)["page"])
	}
}

func TestCreateBookHandlerReturnsFriendlyMessageForInvalidType(t *testing.T) {
	rec, body := postBookRequest(t, `{"name":"Book Shelf","page":80,"type":"magazine","description":"A short description that still meets the minimum length."}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if validationErrors(t, body)["type"] != "Choose a valid book type: fiction, nonfiction, or reference." {
		t.Fatalf("expected friendly type message, got %v", validationErrors(t, body)["type"])
	}
}

func TestCreateBookHandlerReturnsFriendlyMessageForShortDescription(t *testing.T) {
	rec, body := postBookRequest(t, `{"name":"Small Book","page":20,"type":"fiction","description":"Too short"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if validationErrors(t, body)["description"] != "Write a description between 10 and 500 characters." {
		t.Fatalf("expected friendly description message, got %v", validationErrors(t, body)["description"])
	}
}

func TestCreateBookHandlerReturnsFriendlyMessageForMalformedJSON(t *testing.T) {
	rec, body := postBookRequest(t, `{"name":"Broken JSON",`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if body["message"] != "Please send a valid JSON request body." {
		t.Fatalf("expected friendly JSON message, got %v", body["message"])
	}
}

func TestCreateBookHandlerReturnsAllFriendlyMessagesForMultipleErrors(t *testing.T) {
	rec, body := postBookRequest(t, `{"name":"","page":0,"type":"other","description":"short"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	errorsMap := validationErrors(t, body)
	if errorsMap["name"] != "Please enter a book name." {
		t.Fatalf("expected name error, got %v", errorsMap["name"])
	}
	if errorsMap["page"] != "Please enter a page count greater than 0." {
		t.Fatalf("expected page error, got %v", errorsMap["page"])
	}
	if errorsMap["type"] != "Choose a valid book type: fiction, nonfiction, or reference." {
		t.Fatalf("expected type error, got %v", errorsMap["type"])
	}
	if errorsMap["description"] != "Write a description between 10 and 500 characters." {
		t.Fatalf("expected description error, got %v", errorsMap["description"])
	}
}

func TestCreateBookHandlerDoesNotRunNextStepWhenValidationFails(t *testing.T) {
	originalProcessBook := ProcessBook
	defer func() { ProcessBook = originalProcessBook }()

	called := false
	ProcessBook = func(book Book) map[string]interface{} {
		called = true
		return map[string]interface{}{"message": "should not happen"}
	}

	rec, _ := postBookRequest(t, `{"name":"","page":0,"type":"other","description":"short"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if called {
		t.Fatal("expected ProcessBook to stay unused on validation failure")
	}
}
