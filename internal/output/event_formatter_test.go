package output

import (
	"testing"
	"time"

	"dungeon-event-processor/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestFormatEventFormatsIncomingEvents(t *testing.T) {
	tests := []struct {
		name  string
		event domain.OutputEvent
		want  string
	}{
		{
			name: "register",
			event: domain.OutputEvent{
				Timestamp: 14 * time.Hour,
				PlayerID:  1,
				EventType: domain.EventTypeRegister,
			},
			want: "[14:00:00] Player [1] registered",
		},
		{
			name: "cannot continue",
			event: domain.OutputEvent{
				Timestamp: 14*time.Hour + time.Minute,
				PlayerID:  1,
				EventType: domain.EventTypeCannotContinue,
				Extra:     "no mana left",
			},
			want: "[14:01:00] Player [1] cannot continue due to [no mana left]",
		},
		{
			name: "damage",
			event: domain.OutputEvent{
				Timestamp: 14*time.Hour + 2*time.Minute,
				PlayerID:  1,
				EventType: domain.EventTypeDamage,
				Extra:     "60",
			},
			want: "[14:02:00] Player [1] recieved [60] of damage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, FormatEvent(tt.event))
		})
	}
}

func TestFormatEventFormatsOutgoingEvents(t *testing.T) {
	tests := []struct {
		name  string
		event domain.OutputEvent
		want  string
	}{
		{
			name: "disqualified",
			event: domain.OutputEvent{
				Timestamp:    14 * time.Hour,
				PlayerID:     3,
				OutgoingType: domain.OutputEventDisqualified,
			},
			want: "[14:00:00] Player [3] is disqualified",
		},
		{
			name: "dead",
			event: domain.OutputEvent{
				Timestamp:    14 * time.Hour,
				PlayerID:     2,
				OutgoingType: domain.OutputEventDead,
			},
			want: "[14:00:00] Player [2] is dead",
		},
		{
			name: "impossible move",
			event: domain.OutputEvent{
				Timestamp:        14 * time.Hour,
				PlayerID:         2,
				OutgoingType:     domain.OutputEventImpossibleMove,
				RelatedEventType: domain.EventTypePreviousFloor,
			},
			want: "[14:00:00] Player [2] makes imposible move [5]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, FormatEvent(tt.event))
		})
	}
}
