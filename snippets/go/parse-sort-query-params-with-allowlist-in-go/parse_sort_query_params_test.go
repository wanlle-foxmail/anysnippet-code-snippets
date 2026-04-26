package main

import (
	"strings"
	"testing"
)

func TestParseSortQueryParamParsesOneAscendingField(t *testing.T) {
	parsed, err := ParseSortQueryParam("name", []string{"name", "created_at"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(parsed) != 1 || parsed[0].Field != "name" || parsed[0].Desc {
		t.Fatalf("unexpected parsed sort %+v", parsed)
	}
}

func TestParseSortQueryParamParsesMixedDirections(t *testing.T) {
	parsed, err := ParseSortQueryParam("-created_at,name", []string{"name", "created_at"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(parsed) != 2 || !parsed[0].Desc || parsed[0].Field != "created_at" || parsed[1].Desc || parsed[1].Field != "name" {
		t.Fatalf("unexpected parsed sort %+v", parsed)
	}
}

func TestParseSortQueryParamReturnsNilWhenMissing(t *testing.T) {
	parsed, err := ParseSortQueryParam("   ", []string{"name"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if parsed != nil {
		t.Fatalf("expected nil parsed result, got %+v", parsed)
	}
}

func TestParseSortQueryParamRejectsUnknownFields(t *testing.T) {
	_, err := ParseSortQueryParam("priority", []string{"name"})
	if err == nil {
		t.Fatal("expected error for an unknown sort field")
	}
	if !strings.Contains(err.Error(), "not allowed") {
		t.Fatalf("expected allowlist error, got %v", err)
	}
}

func TestParseSortQueryParamRejectsEmptyTerms(t *testing.T) {
	_, err := ParseSortQueryParam("name,,created_at", []string{"name", "created_at"})
	if err == nil {
		t.Fatal("expected error for an empty sort term")
	}
}

func TestParseSortQueryParamRejectsRepeatedFields(t *testing.T) {
	_, err := ParseSortQueryParam("name,-name", []string{"name"})
	if err == nil {
		t.Fatal("expected error for a duplicated sort field")
	}
	if !strings.Contains(err.Error(), "duplicated") {
		t.Fatalf("expected duplicate-field error, got %v", err)
	}
}
