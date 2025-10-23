package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateDashboard(t *testing.T) {
	result := generateDashboard("test", []string{"Startup time (ms)", "Run duration"})

	assert.Contains(t, result, "WHERE MetricName = 'Startup time (ms)'")
	assert.Contains(t, result, "WHERE MetricName = 'Run duration'")
}
