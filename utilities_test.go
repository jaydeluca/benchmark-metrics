package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGenerateTimeframeSlice(t *testing.T) {
	start := "2023-09-20"
	end := "2023-09-23"
	interval := 1
	expected := []string{"2023-09-20", "2023-09-21", "2023-09-22", "2023-09-23"}
	result, _ := generateTimeframeSlice(start, end, interval)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf(fmt.Sprintf("should be equal %s %s", expected, result))
	}
}

func TestConvertDateFormat(t *testing.T) {
	start := "Fri Sep 01 05:16:59 UTC 2023"
	expected := "2023-09-01"
	result := convertDateFormat(start)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf(fmt.Sprintf("should be equal %s %s", expected, result))
	}
}

func TestSplitByMultipleSpaces(t *testing.T) {
	start := `            111.87           129.19           129.52`
	expected := []string{"111.87", "129.19", "129.52"}
	result := splitByMultipleSpaces(start)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf(fmt.Sprintf("should be equal %s %s", expected, result))
	}
}
