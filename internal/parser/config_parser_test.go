package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"dungeon-event-processor/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	path := writeTempConfig(t, `{
		"Floors": 2,
		"Monsters": 3,
		"OpenAt": "14:05:00",
		"Duration": 2
	}`)

	got, err := ParseConfig(path)
	require.NoError(t, err)
	require.Equal(t, domain.DungeonConfig{
		Floors:   2,
		Monsters: 3,
		OpenAt:   14*time.Hour + 5*time.Minute,
		Duration: 2 * time.Hour,
	}, got)
}

func TestParseConfigRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "zero floors",
			body: `{"Floors":0,"Monsters":2,"OpenAt":"14:05:00","Duration":2}`,
		},
		{
			name: "negative monsters",
			body: `{"Floors":2,"Monsters":-1,"OpenAt":"14:05:00","Duration":2}`,
		},
		{
			name: "zero duration",
			body: `{"Floors":2,"Monsters":2,"OpenAt":"14:05:00","Duration":0}`,
		},
		{
			name: "bad open time",
			body: `{"Floors":2,"Monsters":2,"OpenAt":"bad","Duration":2}`,
		},
		{
			name: "bad json",
			body: `{`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseConfig(writeTempConfig(t, tt.body))
			require.Error(t, err)
		})
	}
}

func writeTempConfig(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(path, []byte(body), 0o600)
	require.NoError(t, err)

	return path
}
