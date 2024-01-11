package utils

import (
	"fmt"
	"os"
)

func GetEnv(keys ...string) (map[string]string, error) {
	result := map[string]string{}
	missing := []string{}
	for _, key := range keys {
		value := os.Getenv(key)
		if len(value) == 0 {
			missing = append(missing, key)
		} else {
			result[key] = value
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing environment variables: %v", missing)
	}
	return result, nil
}
