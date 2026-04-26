package main

import (
	"testing"
)

func TestPickAllowedFieldsFromPatchBodySelectsAllowedFields(t *testing.T) {
	body := map[string]interface{}{"display_name": "Ada", "role": "admin"}
	picked, err := PickAllowedFieldsFromPatchBody(body, []string{"display_name"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(picked) != 1 || picked["display_name"] != "Ada" {
		t.Fatalf("unexpected picked fields %v", picked)
	}
}

func TestPickAllowedFieldsFromPatchBodyPreservesExplicitValues(t *testing.T) {
	body := map[string]interface{}{"active": false, "nickname": nil, "attempts": float64(0)}
	picked, err := PickAllowedFieldsFromPatchBody(body, []string{"active", "nickname", "attempts"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if picked["active"] != false || picked["nickname"] != nil || picked["attempts"] != float64(0) {
		t.Fatalf("unexpected picked values %v", picked)
	}
}

func TestPickAllowedFieldsFromPatchBodyIgnoresUnknownFields(t *testing.T) {
	body := map[string]interface{}{"display_name": "Ada", "role": "admin"}
	picked, err := PickAllowedFieldsFromPatchBody(body, []string{"display_name"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if _, ok := picked["role"]; ok {
		t.Fatalf("expected unknown field role to be ignored, got %v", picked)
	}
}

func TestPickAllowedFieldsFromPatchBodyReturnsEmptyMapForNilBody(t *testing.T) {
	picked, err := PickAllowedFieldsFromPatchBody(nil, []string{"display_name"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(picked) != 0 {
		t.Fatalf("expected an empty map, got %v", picked)
	}
}

func TestPickAllowedFieldsFromPatchBodyRejectsNonMappingInput(t *testing.T) {
	_, err := PickAllowedFieldsFromPatchBody([]string{"display_name"}, []string{"display_name"})
	if err == nil {
		t.Fatal("expected error for a non-map patch body")
	}
}

func TestPickAllowedFieldsFromPatchBodyRejectsBlankAllowedFieldNames(t *testing.T) {
	_, err := PickAllowedFieldsFromPatchBody(map[string]interface{}{}, []string{"display_name", "  "})
	if err == nil {
		t.Fatal("expected error for a blank allowed field name")
	}
}
