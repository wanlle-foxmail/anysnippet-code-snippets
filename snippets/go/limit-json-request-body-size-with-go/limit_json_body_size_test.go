package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONBodyWithLimitAcceptsBodyUnderLimit(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"Ada"}`))
	recorder := httptest.NewRecorder()

	var body CreateBookRequest
	err := DecodeJSONBodyWithLimit(recorder, request, 32, &body)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if body.Title != "Ada" {
		t.Fatalf("expected title Ada, got %+v", body)
	}
}

func TestDecodeJSONBodyWithLimitAcceptsBodyAtExactLimit(t *testing.T) {
	payload := `{"title":"1234567890"}`
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(payload))
	recorder := httptest.NewRecorder()

	var body CreateBookRequest
	err := DecodeJSONBodyWithLimit(recorder, request, int64(len(payload)), &body)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDecodeJSONBodyWithLimitRejectsOversizedBodies(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"this title is too large"}`))
	recorder := httptest.NewRecorder()

	var body CreateBookRequest
	err := DecodeJSONBodyWithLimit(recorder, request, 10, &body)
	if err == nil {
		t.Fatal("expected an error for an oversized body")
	}
	var maxBytesErr *http.MaxBytesError
	if !errors.As(err, &maxBytesErr) {
		t.Fatalf("expected MaxBytesError, got %v", err)
	}
}

func TestDecodeJSONBodyWithLimitRejectsMalformedJSONUnderLimit(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":`))
	recorder := httptest.NewRecorder()

	var body CreateBookRequest
	err := DecodeJSONBodyWithLimit(recorder, request, 64, &body)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestDecodeJSONBodyWithLimitRejectsMultipleJSONValues(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"Ada"} {"title":"Grace"}`))
	recorder := httptest.NewRecorder()

	var body CreateBookRequest
	err := DecodeJSONBodyWithLimit(recorder, request, 64, &body)
	if err == nil {
		t.Fatal("expected error for multiple JSON values")
	}
	if !strings.Contains(err.Error(), "exactly one JSON value") {
		t.Fatalf("expected single-value error, got %v", err)
	}
}

func TestDecodeJSONBodyWithLimitRejectsInvalidInput(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"Ada"}`))
	recorder := httptest.NewRecorder()

	if err := DecodeJSONBodyWithLimit(nil, request, 64, &CreateBookRequest{}); err == nil {
		t.Fatal("expected error for nil writer")
	}
	if err := DecodeJSONBodyWithLimit(recorder, nil, 64, &CreateBookRequest{}); err == nil {
		t.Fatal("expected error for nil request")
	}
	if err := DecodeJSONBodyWithLimit(recorder, request, 64, nil); err == nil {
		t.Fatal("expected error for nil target")
	}
	if err := DecodeJSONBodyWithLimit(recorder, request, 0, &CreateBookRequest{}); err == nil {
		t.Fatal("expected error for invalid max bytes")
	}
}
