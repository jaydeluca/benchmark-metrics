package main

import (
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"testing"
	"time"
)

func TestConvertReportToDataPoint(t *testing.T) {
	layout := "2006-01-02"
	dateString := "2023-09-01"
	date, _ := time.Parse(layout, dateString)
	expected := metricdata.DataPoint[float64]{
		Attributes: attribute.NewSet(attribute.String("entity", "none")),
		StartTime:  date,
		Time:       time.Now(),
		Value:      0.51,
	}
	result := generateDataPoint("none", date, 0.51)
	assert.Equal(t, expected.StartTime, result.StartTime)
	assert.Equal(t, expected.Value, result.Value)
	assert.True(t, expected.Attributes.HasValue("entity"))
	assert.True(t, result.Attributes.HasValue("entity"))
}
