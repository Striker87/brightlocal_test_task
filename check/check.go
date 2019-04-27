package check

import "strings"

func KeyExists(key string, storage map[string]string) bool {
	_, ok := storage[key]
	return ok
}

func Key(key string) bool {
	if strings.TrimSpace(key) == "" || len(key) > 16 {
		return true
	}
	return false
}

func Value(value string) bool {
	return len(value) > 512
}

func StorageSize(storage map[string]string) bool {
	return len(storage) > 1024
}
