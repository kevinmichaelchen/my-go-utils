package utils

import (
	"os"
	"log"
	"strconv"
)

// EnvOrBool returns the environment variable for the given key,
// or the default value if no environment variable is found.
func EnvOrBool(key string, defaultVal bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("No env var found for %s. Using default value: %t\n", key, defaultVal)
		return defaultVal
	}
	if b, err := strconv.ParseBool(v); err != nil {
		log.Printf("Env var for %s is not a proper boolean\n", key)
		panic(err)
	} else {
		return b
	}
}

// EnvOrInt returns the environment variable for the given key as an int,
// or the default value if no environment variable is found.
func EnvOrInt(key string, defaultVal int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("No env var found for %s. Using default value: %d\n", key, defaultVal)
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("Env var %s needs to be an int...\n", key)
	}
	return i
}

// EnvOrString returns the environment variable for the given key,
// or the default value if no environment variable is found.
func EnvOrString(key, defaultVal string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("No env var found for %s. Using default value: %s\n", key, defaultVal)
		return defaultVal
	}
	return v
}

// MustEnv returns the environment variable for the given key, or exits if no such variable is found.
func MustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("No env var for: %s", key)
	}
	return v
}