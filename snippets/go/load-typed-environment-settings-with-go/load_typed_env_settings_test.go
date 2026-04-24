package main

import (
	"strings"
	"testing"
)

func TestLoadTypedEnvSettingsLoadsRequiredAndOptionalValues(t *testing.T) {
	settings, err := LoadTypedEnvSettings(map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
		"DEBUG":   "true",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := AppSettings{AppEnv: "production", Port: 8080, Debug: true}
	if settings != want {
		t.Fatalf("got %+v; want %+v", settings, want)
	}
}

func TestLoadTypedEnvSettingsUsesDefaultsForMissingOptionalValues(t *testing.T) {
	settings, err := LoadTypedEnvSettings(map[string]string{
		"APP_ENV": "development",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := AppSettings{AppEnv: "development", Port: 8000, Debug: false}
	if settings != want {
		t.Fatalf("got %+v; want %+v", settings, want)
	}
}

func TestLoadTypedEnvSettingsReturnsErrorWhenRequiredValueIsMissing(t *testing.T) {
	_, err := LoadTypedEnvSettings(map[string]string{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "APP_ENV") {
		t.Fatalf("expected APP_ENV error, got %v", err)
	}
}

func TestLoadTypedEnvSettingsReturnsErrorForInvalidPort(t *testing.T) {
	_, err := LoadTypedEnvSettings(map[string]string{
		"APP_ENV": "test",
		"PORT":    "abc",
	})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "PORT") {
		t.Fatalf("expected PORT error, got %v", err)
	}
}

func TestLoadTypedEnvSettingsReturnsErrorForInvalidDebugValue(t *testing.T) {
	_, err := LoadTypedEnvSettings(map[string]string{
		"APP_ENV": "test",
		"DEBUG":   "sometimes",
	})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "DEBUG") {
		t.Fatalf("expected DEBUG error, got %v", err)
	}
}

func TestLoadTypedEnvSettingsReadsFromOSEnvironmentWhenMapIsNil(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	t.Setenv("PORT", "9000")
	t.Setenv("DEBUG", "0")

	settings, err := LoadTypedEnvSettings(nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := AppSettings{AppEnv: "staging", Port: 9000, Debug: false}
	if settings != want {
		t.Fatalf("got %+v; want %+v", settings, want)
	}
}