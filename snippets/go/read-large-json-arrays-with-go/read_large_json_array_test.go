package main

import (
	"errors"
	"strings"
	"testing"
)

func TestReadJSONArrayItemsReadsItemsInOrder(t *testing.T) {
	reader := strings.NewReader(`[{"id":1},{"id":2},{"id":3}]`)
	var ids []int

	err := ReadJSONArrayItems(reader, func(item Event) error {
		ids = append(ids, item.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(ids) != 3 || ids[0] != 1 || ids[1] != 2 || ids[2] != 3 {
		t.Fatalf("unexpected ids %v", ids)
	}
}

func TestReadJSONArrayItemsKeepsNestedObjectsAndLists(t *testing.T) {
	reader := strings.NewReader(`[{"id":1,"tags":["a","b"],"meta":{"source":"api"}}]`)
	var item Event

	err := ReadJSONArrayItems(reader, func(decoded Event) error {
		item = decoded
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(item.Tags) != 2 || item.Tags[0] != "a" || item.Meta["source"] != "api" {
		t.Fatalf("unexpected item %+v", item)
	}
}

func TestReadJSONArrayItemsHandlesEmptyArrays(t *testing.T) {
	reader := strings.NewReader(`[]`)
	called := false

	err := ReadJSONArrayItems(reader, func(item Event) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if called {
		t.Fatal("expected the handler to stay unused for an empty array")
	}
}

func TestReadJSONArrayItemsRejectsMalformedArrays(t *testing.T) {
	reader := strings.NewReader(`[{"id":1}`)
	err := ReadJSONArrayItems(reader, func(Event) error { return nil })
	if err == nil {
		t.Fatal("expected error for malformed JSON array")
	}
}

func TestReadJSONArrayItemsRejectsNonArrayInput(t *testing.T) {
	reader := strings.NewReader(`{"id":1}`)
	err := ReadJSONArrayItems(reader, func(Event) error { return nil })
	if err == nil {
		t.Fatal("expected error for non-array input")
	}
}

func TestReadJSONArrayItemsRejectsTrailingContent(t *testing.T) {
	reader := strings.NewReader(`[{"id":1}] {"extra":"data"}`)
	err := ReadJSONArrayItems(reader, func(Event) error { return nil })
	if err == nil {
		t.Fatal("expected error for trailing content after the array")
	}
	if !strings.Contains(err.Error(), "exactly one JSON array") {
		t.Fatalf("expected single-array error, got %v", err)
	}
}

func TestReadJSONArrayItemsPropagatesHandlerErrors(t *testing.T) {
	reader := strings.NewReader(`[{"id":1}]`)
	wantErr := errors.New("write failed")
	err := ReadJSONArrayItems(reader, func(Event) error { return wantErr })
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}

func TestReadJSONArrayItemsRejectsInvalidInput(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		if err := ReadJSONArrayItems[Event](nil, func(Event) error { return nil }); err == nil {
			t.Fatal("expected error for a nil reader")
		}
	})

	t.Run("nil handler", func(t *testing.T) {
		if err := ReadJSONArrayItems[Event](strings.NewReader(`[]`), nil); err == nil {
			t.Fatal("expected error for a nil handler")
		}
	})
}
