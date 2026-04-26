package main

import (
	"net/url"
	"strings"
	"testing"
)

func TestParsePaginationQueryParamsParsesValidValues(t *testing.T) {
	values := url.Values{"page": {"3"}, "page_size": {"20"}}

	pagination, err := ParsePaginationQueryParams(values, 1, 25, 50)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if pagination.Page != 3 || pagination.PageSize != 20 || pagination.Offset != 40 {
		t.Fatalf("unexpected pagination %+v", pagination)
	}
}

func TestParsePaginationQueryParamsUsesDefaultsWhenMissing(t *testing.T) {
	pagination, err := ParsePaginationQueryParams(url.Values{}, 1, 25, 50)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if pagination.Page != 1 || pagination.PageSize != 25 || pagination.Offset != 0 {
		t.Fatalf("unexpected pagination %+v", pagination)
	}
}

func TestParsePaginationQueryParamsCapsLargePageSizes(t *testing.T) {
	values := url.Values{"page_size": {"500"}}

	pagination, err := ParsePaginationQueryParams(values, 1, 25, 100)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if pagination.PageSize != 100 {
		t.Fatalf("expected capped page size 100, got %d", pagination.PageSize)
	}
}

func TestParsePaginationQueryParamsRejectsInvalidIntegers(t *testing.T) {
	_, err := ParsePaginationQueryParams(url.Values{"page": {"nope"}}, 1, 25, 50)
	if err == nil {
		t.Fatal("expected error for invalid page")
	}
	if !strings.Contains(err.Error(), "integer") {
		t.Fatalf("expected integer error, got %v", err)
	}
}

func TestParsePaginationQueryParamsRejectsNonPositiveValues(t *testing.T) {
	_, err := ParsePaginationQueryParams(url.Values{"page_size": {"0"}}, 1, 25, 50)
	if err == nil {
		t.Fatal("expected error for non-positive page size")
	}
	if !strings.Contains(err.Error(), "greater than 0") {
		t.Fatalf("expected positive-value error, got %v", err)
	}
}

func TestParsePaginationQueryParamsRejectsInvalidDefaults(t *testing.T) {
	t.Run("default page", func(t *testing.T) {
		if _, err := ParsePaginationQueryParams(url.Values{}, 0, 25, 50); err == nil {
			t.Fatal("expected error for invalid default page")
		}
	})

	t.Run("default page size", func(t *testing.T) {
		if _, err := ParsePaginationQueryParams(url.Values{}, 1, 0, 50); err == nil {
			t.Fatal("expected error for invalid default page size")
		}
	})

	t.Run("max page size", func(t *testing.T) {
		if _, err := ParsePaginationQueryParams(url.Values{}, 1, 25, 0); err == nil {
			t.Fatal("expected error for invalid max page size")
		}
	})
}

func TestParsePaginationQueryParamsTrimsSurroundingWhitespace(t *testing.T) {
	values := url.Values{"page": {" 4 "}, "page_size": {" 10 "}}

	pagination, err := ParsePaginationQueryParams(values, 1, 25, 50)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if pagination.Page != 4 || pagination.PageSize != 10 || pagination.Offset != 30 {
		t.Fatalf("unexpected pagination %+v", pagination)
	}
}
