package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleStrLowerTrim() {
	StrLowerTrim("   HeLLo  ") // hello
}

func TestStrLowerTrim(t *testing.T) {
	assert.Equal(t, StrLowerTrim("   HeLLo  "), "hello")
	assert.NotEqual(t, StrLowerTrim("   HeLLo!  "), "HELLO!")
}

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

func ExampleStrIsIn() {
	listOfStrings := []string{"hello", "world"}
	if StrIsIn("hello", listOfStrings) {
		// ... do your stuff
	}
}

func TestStrIsIn(t *testing.T) {
	assert.True(t, StrIsIn("hello", []string{"hello", "world"}))
	assert.False(t, StrIsIn("hello", []string{"Hello", "world"}))
}

func ExampleStrIsInEqualFold() {
	listOfStrings := []string{"hello", "world"}
	if StrIsInEqualFold("hello", listOfStrings) {
		// ... do your stuff
	}

	listOfStringsUpper := []string{"HELLO", "WORLD"}
	if StrIsInEqualFold("hello", listOfStringsUpper) {
		// ... this will be true too!
	}
}

func TestStrIsInEqualFold(t *testing.T) {
	assert.True(t, StrIsInEqualFold("hello", []string{"hello", "world"}))
	assert.True(t, StrIsInEqualFold("hello", []string{"Hello", "world"}))
	assert.False(t, StrIsInEqualFold("hello", []string{"Python", "world"}))
}

func ExampleStrUpperTrim() {
	StrUpperTrim("   HeLLo  ") // HELLO
}

func TestStrUpperTrim(t *testing.T) {
	assert.Equal(t, StrUpperTrim("   HeLLo  "), "HELLO")
	assert.NotEqual(t, StrUpperTrim("   HeLLo!  "), "hello!")
}
