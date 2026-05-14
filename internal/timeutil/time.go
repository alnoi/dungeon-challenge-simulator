package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseClock(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "[]")

	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("clock must use HH:MM:SS format: %q", value)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("parse hours: %w", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("parse minutes: %w", err)
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("parse seconds: %w", err)
	}

	if hours < 0 || minutes < 0 || minutes > 59 || seconds < 0 || seconds > 59 {
		return 0, fmt.Errorf("invalid clock value: %q", value)
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}

func FormatDuration(duration time.Duration) string {
	if duration < 0 {
		duration = 0
	}

	totalSeconds := int64(duration / time.Second)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func FormatTimestamp(timestamp time.Duration) string {
	return "[" + FormatDuration(timestamp) + "]"
}
