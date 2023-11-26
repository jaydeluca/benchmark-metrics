package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCacheAdd(t *testing.T) {
	cache := NewSingleFileCache("test-cache")
	err := cache.AddToCache("test", "value")
	if err != nil {
		t.Failed()
	}

	value, err := cache.RetrieveValue("test")
	if err != nil {
		t.Failed()
	}

	assert.Equal(t, "value", value)
	empty, _ := cache.RetrieveValue("not existent")
	assert.Empty(t, empty)
	err = cache.DeleteCache()
	if err != nil {
		t.Failed()
	}
}
