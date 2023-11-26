package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateDashboard(t *testing.T) {
	result := generateDashboard("test", []string{"cpu", "memory"})

	assert.Contains(t, result, "WHERE MetricName = 'cpu'")
	assert.Contains(t, result, "WHERE MetricName = 'memory'")
}
