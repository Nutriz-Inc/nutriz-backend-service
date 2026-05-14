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

func IsFutureDate(dateStr string) bool {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return false
	}
	return date.After(time.Now())
}

func StringToTime(dateStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateStr)
}
