package utils

import (
	"os"
)

// Returns true if path or file exists. Otherwise false
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
