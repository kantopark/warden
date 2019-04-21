package utils

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamplePathExists() {
	PathExists("path/to/file/exists")         // true
	PathExists("path/to/file/does/not/exist") // false
}

func TestPathExists(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	assert.True(t, PathExists(file))

	dir := filepath.Dir(file)
	assert.True(t, PathExists(dir))

	assert.False(t, PathExists(filepath.Join(file, "CONFIRM_DOES_NOT_EXISTS.pdf")))
}
