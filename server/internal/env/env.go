package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

func GetInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return intValue
}

func GetBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}

	return boolValue
}

func GetDuration(key string, unit time.Duration, defaultValue int) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return time.Duration(defaultValue) * unit
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return time.Duration(intValue) * unit
}
