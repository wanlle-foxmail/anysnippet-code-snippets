package main

import (
	"strings"
	"testing"
)

func TestDecodeStrictJSONAcceptsKnownFields(t *testing.T) {
	var request CreateUserRequest
	err := DecodeStrictJSON(
		strings.NewReader(`{"display_name":"Ada","email":"ada@example.com"}`),
		&request,
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if request.DisplayName != "Ada" || request.Email != "ada@example.com" {
		t.Fatalf("got %+v; want populated request", request)
	}
}

func TestDecodeStrictJSONRejectsUnknownFields(t *testing.T) {
	var request CreateUserRequest
	err := DecodeStrictJSON(
		strings.NewReader(`{"display_name":"Ada","role":"admin"}`),
		&request,
	)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("expected unknown field error, got %v", err)
	}
}

func TestDecodeStrictJSONRejectsMalformedJSON(t *testing.T) {
	var request CreateUserRequest
	err := DecodeStrictJSON(strings.NewReader(`{"display_name":`), &request)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestDecodeStrictJSONRejectsMultipleJSONValues(t *testing.T) {
	var request CreateUserRequest
	err := DecodeStrictJSON(
		strings.NewReader(`{"display_name":"Ada"} {"email":"ada@example.com"}`),
		&request,
	)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "exactly one JSON value") {
		t.Fatalf("expected multiple values error, got %v", err)
	}
}

func TestDecodeStrictJSONRejectsNilReader(t *testing.T) {
	var request CreateUserRequest
	err := DecodeStrictJSON(nil, &request)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "reader") {
		t.Fatalf("expected reader error, got %v", err)
	}
}

func TestDecodeStrictJSONRejectsNilTarget(t *testing.T) {
	err := DecodeStrictJSON(strings.NewReader(`{"display_name":"Ada"}`), nil)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "target") {
		t.Fatalf("expected target error, got %v", err)
	}
}