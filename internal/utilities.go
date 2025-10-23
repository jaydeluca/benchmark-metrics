package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	layout = "2006-01-02"
)

func generateTimeframeToToday(start string, interval int) ([]string, error) {
	currentTime := time.Now()
	return generateTimeframeSlice(start, currentTime.Format(layout), interval)
}

func generateTimeframeSlice(start, end string, interval int) ([]string, error) {
	startDate, err := time.Parse(layout, start)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start date: %v", err)
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end date: %v", err)
	}

	var dateList []string

	// Increment the start date by the interval until it reaches or exceeds the end date
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		dateList = append(dateList, currentDate.Format(layout))
		currentDate = currentDate.AddDate(0, 0, interval)
	}

	return dateList, nil
}

func convertDateFormat(inputDate string) time.Time {
	layout := "Mon Jan 02 15:04:05 MST 2006"
	parsedTime, err := time.Parse(layout, inputDate)
	if err != nil {
		panic(err)
	}

	return parsedTime
}

func splitByMultipleSpaces(input string) []string {
	// Split by one or more whitespaces
	regex := regexp.MustCompile(`\s+`)
	splitValues := regex.Split(input, -1)

	// Filter out empty strings and trim spaces
	var cleanedValues []string
	for _, value := range splitValues {
		trimmedValue := strings.TrimSpace(value)
		if trimmedValue != "" {
			cleanedValues = append(cleanedValues, trimmedValue)
		}
	}
	return cleanedValues
}

// parseTimestampToSeconds converts a timestamp in HH:MM:SS format to total minutes
// Returns -1 if the timestamp cannot be parsed
func parseTimestampToSeconds(timestamp string) float64 {
	parts := strings.Split(strings.TrimSpace(timestamp), ":")
	if len(parts) != 3 {
		return -1
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return -1
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return -1
	}

	// Convert to total minutes: hours * 60 + minutes + seconds / 60
	totalMinutes := float64(hours*60) + float64(minutes) + float64(seconds)/60.0
	return totalMinutes
}
