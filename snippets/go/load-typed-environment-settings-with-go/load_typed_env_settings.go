package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type AppSettings struct {
	AppEnv string
	Port   int
	Debug  bool
}

func LoadTypedEnvSettings(env map[string]string) (AppSettings, error) {
	// Flow:
	//   read APP_ENV, PORT, and DEBUG from the provided map or OS environment
	//      |
	//      +-> parse text, int, and bool values -> return AppSettings
	//      `-> missing or invalid values -> return an error
	appEnv, err := readRequiredText(env, "APP_ENV")
	if err != nil {
		return AppSettings{}, err
	}

	port, err := readInt(env, "PORT", 8000)
	if err != nil {
		return AppSettings{}, err
	}

	debug, err := readBool(env, "DEBUG", false)
	if err != nil {
		return AppSettings{}, err
	}

	return AppSettings{AppEnv: appEnv, Port: port, Debug: debug}, nil
}

func readRequiredText(env map[string]string, key string) (string, error) {
	value, ok := lookupValue(env, key)
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return strings.TrimSpace(value), nil
}

func readInt(env map[string]string, key string, defaultValue int) (int, error) {
	value, ok := lookupValue(env, key)
	if !ok {
		return defaultValue, nil
	}

	parsedValue, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer", key)
	}
	return parsedValue, nil
}

func readBool(env map[string]string, key string, defaultValue bool) (bool, error) {
	value, ok := lookupValue(env, key)
	if !ok {
		return defaultValue, nil
	}

	normalizedValue := strings.ToLower(strings.TrimSpace(value))
	if normalizedValue == "1" || normalizedValue == "true" {
		return true, nil
	}
	if normalizedValue == "0" || normalizedValue == "false" {
		return false, nil
	}

	return false, fmt.Errorf("%s must be one of: true, false, 1, 0", key)
}

func lookupValue(env map[string]string, key string) (string, bool) {
	if env != nil {
		value, ok := env[key]
		return value, ok
	}
	return os.LookupEnv(key)
}

func main() {
	settings, err := LoadTypedEnvSettings(map[string]string{
		"APP_ENV": "development",
		"PORT":    "8000",
		"DEBUG":   "false",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("settings=%+v", settings)
}