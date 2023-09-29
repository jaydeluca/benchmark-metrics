package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCacheAdd(t *testing.T) {
	cache := NewSingleFileCache("test-cache")
	cache.AddToCache("test", "value")
	value, _ := cache.RetrieveValue("test")
	assert.Equal(t, "value", value)
	empty, _ := cache.RetrieveValue("not existent")
	assert.Empty(t, empty)
	cache.DeleteCache()
}
