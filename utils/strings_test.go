package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleStrIn() {
	StrIn("Hello", nil, "hello")
	StrIn("Hello", &StrOption{IgnoreCase: false}, "hello")
}

func TestStrIn(t *testing.T) {
	assert.True(t, StrIn("Hello", nil, "hello"))
	assert.False(t, StrIn("Hello", &StrOption{IgnoreCase: false}, "hello"))
	assert.False(t, StrIn("Hello", nil))
}

func ExampleStrEquals() {
	StrEquals("hello", "Hello", nil)
	StrEquals("hello", "Hello", &StrOption{IgnoreCase: false})
}

func TestStrEquals(t *testing.T) {
	assert.False(t, StrEquals("Hello", "hello", &StrOption{IgnoreCase: false}))
	assert.True(t, StrEquals("Hello", "hello", &StrOption{IgnoreCase: true}))
	assert.True(t, StrEquals("Hello", "hello", nil))
}

func TestStrIsEmptyOrWhitespace(t *testing.T) {
	assert.False(t, StrIsEmptyOrWhitespace("a "))
	assert.True(t, StrIsEmptyOrWhitespace(""))
	assert.True(t, StrIsEmptyOrWhitespace("  \n "))
}
