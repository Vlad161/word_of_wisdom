package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

func MustString(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required ENV %q is not set", key)
	}
	if v == "" {
		log.Fatalf("required ENV %q is empty", key)
	}
	return v
}

func MustBool(key string) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required ENV %q is not set", key)
	}
	if v == "true" || v == "1" {
		return true
	}
	return false
}

func MustInt(key string) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required ENV %q is not set", key)
	}
	if v == "" {
		log.Fatalf("required ENV %q is empty", key)
	}
	res, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		log.Fatalf("required ENV %q must be a number but it's %q", key, v)
	}
	return int(res)
}

func MustFloat(key string) float64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required ENV %q is not set", key)
	}
	if v == "" {
		log.Fatalf("required ENV %q is empty", key)
	}
	res, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Fatalf("required ENV %q must be a float but it's %q", key, v)
	}
	return res
}

func MustDuration(key string) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required ENV %q is not set", key)
	}

	res, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("required ENV %q must be a parsable duration but it's %q: %v", key, v, err)
	}
	return res
}
