package deploy

import (
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

var rm *routeMap

func init() {
	rm = newRouteMap()
	rm.Set("hello", "world")
}

func TestRouteMap_Get(t *testing.T) {
	assert.Equal(t, rm.Get("hello"), "world")
	assert.Equal(t, rm.Get("no-exist"), "")
}

func TestRouteMap_Set(t *testing.T) {
	testName, value := "set-test", "set-test-value"
	rm.Set(testName, value)
	assert.Equal(t, rm.Get(testName), value)
}

func TestRouteMap_Delete(t *testing.T) {
	testName, value := "del-test", "del-test-value"
	rm.Set(testName, value)
	assert.Equal(t, rm.Get(testName), value)
	rm.Delete("del-test")
	assert.Equal(t, rm.Get(testName), "")
}
