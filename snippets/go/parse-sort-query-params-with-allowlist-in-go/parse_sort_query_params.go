package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type SortField struct {
	Field string
	Desc  bool
}

// Flow: split the raw sort string -> validate each token against the allowlist -> return ordered asc or desc fields.
func ParseSortQueryParam(raw string, allowlist []string) ([]SortField, error) {
	allowedFields := make(map[string]struct{}, len(allowlist))
	for _, field := range allowlist {
		trimmed := strings.TrimSpace(field)
		if trimmed == "" {
			return nil, errors.New("allowlist fields must not be blank")
		}
		allowedFields[trimmed] = struct{}{}
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	terms := strings.Split(raw, ",")
	parsed := make([]SortField, 0, len(terms))
	seen := make(map[string]struct{}, len(terms))

	for _, term := range terms {
		trimmed := strings.TrimSpace(term)
		if trimmed == "" {
			return nil, errors.New("sort terms must not be empty")
		}

		sortField := SortField{Field: trimmed}
		if strings.HasPrefix(trimmed, "-") {
			sortField.Field = strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
			sortField.Desc = true
		}
		if sortField.Field == "" {
			return nil, errors.New("sort field is required")
		}
		if _, ok := allowedFields[sortField.Field]; !ok {
			return nil, fmt.Errorf("sort field %q is not allowed", sortField.Field)
		}
		if _, ok := seen[sortField.Field]; ok {
			return nil, fmt.Errorf("sort field %q is duplicated", sortField.Field)
		}

		seen[sortField.Field] = struct{}{}
		parsed = append(parsed, sortField)
	}

	return parsed, nil
}

func main() {
	sortFields, err := ParseSortQueryParam("-created_at,name", []string{"created_at", "name", "priority"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sort=%+v", sortFields)
}
