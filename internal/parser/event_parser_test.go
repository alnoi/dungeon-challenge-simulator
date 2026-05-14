package parser

import (
	"strings"
	"testing"
	"time"

	"dungeon-event-processor/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestParseEvent(t *testing.T) {
	tests := []struct {
		name string
		line string
		want domain.Event
	}{
		{
			name: "without extra param",
			line: "[14:10:00] 2 3",
			want: domain.Event{
				Timestamp: 14*time.Hour + 10*time.Minute,
				PlayerID:  2,
				Type:      domain.EventTypeKillMonster,
			},
		},
		{
			name: "with multi word extra param",
			line: "[14:10:00] 2 9 no mana left",
			want: domain.Event{
				Timestamp: 14*time.Hour + 10*time.Minute,
				PlayerID:  2,
				Type:      domain.EventTypeCannotContinue,
				Extra:     "no mana left",
			},
		},
		{
			name: "with numeric extra param",
			line: "[14:10:00] 2 11 60",
			want: domain.Event{
				Timestamp: 14*time.Hour + 10*time.Minute,
				PlayerID:  2,
				Type:      domain.EventTypeDamage,
				Extra:     "60",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEvent(tt.line)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseEventRejectsInvalidLines(t *testing.T) {
	tests := []string{
		"",
		"[14:10:00] bad 3",
		"[14:10:00] 2 bad",
		"[14:10:00] 2 99",
		"[14:99:00] 2 3",
	}

	for _, line := range tests {
		t.Run(line, func(t *testing.T) {
			_, err := ParseEvent(line)
			require.Error(t, err)
		})
	}
}

func TestParseEvents(t *testing.T) {
	input := strings.NewReader(`
[14:00:00] 1 1

[14:05:00] 1 2
`)

	events, err := ParseEvents(input)
	require.NoError(t, err)
	require.Equal(t, []domain.Event{
		{
			Timestamp: 14 * time.Hour,
			PlayerID:  1,
			Type:      domain.EventTypeRegister,
		},
		{
			Timestamp: 14*time.Hour + 5*time.Minute,
			PlayerID:  1,
			Type:      domain.EventTypeEnterDungeon,
		},
	}, events)
}
