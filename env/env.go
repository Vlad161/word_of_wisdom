package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return v
}

func GetBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	if v == "true" || v == "1" {
		return true
	}
	return false
}

func GetInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	res, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return fallback
	}
	return int(res)
}

func GetFloat(key string, fallback float64) float64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	res, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return res
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	res, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return res
}
