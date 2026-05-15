package utils

import (
	"strconv"
	"time"
)

func IntPtr(v int) *int {
	return &v
}

func Float64Ptr(v float64) *float64 {
	return &v
}

func StringPtr(v string) *string {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}

func StringToFloat64(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	return parsed
}

func StringToDate(value string) (*time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func StringToTime(value string) (*time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func IsFutureDate(dateStr string) bool {
	for _, layout := range []string{"2006-01-02", "2006-01-02 15:04:05", time.RFC3339} {
		date, err := time.Parse(layout, dateStr)
		if err == nil {
			return date.After(time.Now())
        }
	}
	return false
}