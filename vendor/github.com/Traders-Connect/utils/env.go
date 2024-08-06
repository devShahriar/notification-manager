package utils

import (
	"os"
	"strconv"
)

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvOrUInt64(key string, defaultVal uint64) (uint64, error) {
	if val, ok := os.LookupEnv(key); ok {
		return strconv.ParseUint(val, 10, 64)
	}
	return defaultVal, nil
}

func LookupEnvOrInt64(key string, defaultVal int64) (int64, error) {
	if val, ok := os.LookupEnv(key); ok {
		return strconv.ParseInt(val, 10, 64)
	}
	return defaultVal, nil
}

func LookupEnvOrBool(key string, defaultVal bool) (bool, error) {
	if val, ok := os.LookupEnv(key); ok {
		return strconv.ParseBool(val)
	}
	return defaultVal, nil
}
