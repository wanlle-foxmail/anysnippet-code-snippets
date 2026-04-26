package main

import (
	"errors"
	"log"
	"strings"
)

// Flow: inspect the PATCH body -> keep only allowlisted top-level keys -> preserve explicit values and ignore unknown fields.
func PickAllowedFieldsFromPatchBody(body interface{}, allowedFields []string) (map[string]interface{}, error) {
	allowed := make(map[string]struct{}, len(allowedFields))
	for _, field := range allowedFields {
		trimmed := strings.TrimSpace(field)
		if trimmed == "" {
			return nil, errors.New("allowed fields must not be blank")
		}
		allowed[trimmed] = struct{}{}
	}

	if body == nil {
		return map[string]interface{}{}, nil
	}

	patchBody, ok := body.(map[string]interface{})
	if !ok {
		return nil, errors.New("patch body must be a map[string]interface{}")
	}

	picked := make(map[string]interface{})
	for key, value := range patchBody {
		if _, ok := allowed[key]; ok {
			picked[key] = value
		}
	}

	return picked, nil
}

func main() {
	body := map[string]interface{}{
		"display_name": "Ada",
		"role":         "admin",
		"nickname":     nil,
	}

	picked, err := PickAllowedFieldsFromPatchBody(body, []string{"display_name", "nickname"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("picked=%v", picked)
}
