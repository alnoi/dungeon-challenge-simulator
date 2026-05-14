package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseClock(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Duration
	}{
		{
			name:  "plain clock",
			input: "14:05:09",
			want:  14*time.Hour + 5*time.Minute + 9*time.Second,
		},
		{
			name:  "timestamp brackets",
			input: "[00:01:02]",
			want:  time.Minute + 2*time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseClock(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseClockRejectsInvalidValues(t *testing.T) {
	tests := []string{
		"14:05",
		"14:60:00",
		"14:00:60",
		"bad",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseClock(input)
			require.Error(t, err)
		})
	}
}

func TestFormatDuration(t *testing.T) {
	require.Equal(t, "14:05:09", FormatDuration(14*time.Hour+5*time.Minute+9*time.Second))
	require.Equal(t, "00:00:00", FormatDuration(-time.Second))
}

func TestFormatTimestamp(t *testing.T) {
	require.Equal(t, "[01:02:03]", FormatTimestamp(time.Hour+2*time.Minute+3*time.Second))
}
