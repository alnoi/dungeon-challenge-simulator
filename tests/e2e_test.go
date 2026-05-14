package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dungeon-event-processor/internal/engine"
	"dungeon-event-processor/internal/output"
	"dungeon-event-processor/internal/parser"

	"github.com/stretchr/testify/require"
)

func TestDungeonFlowEndToEnd(t *testing.T) {
	configPath := writeTempFile(t, "config.json", `{
		"Floors": 3,
		"Monsters": 2,
		"OpenAt": "10:00:00",
		"Duration": 1
	}`)

	eventsPath := writeTempFile(t, "events", `[09:59:00] 1 1
[10:00:00] 1 2
[10:01:00] 2 1
[10:02:00] 1 3
[10:02:00] 2 2
[10:02:30] 3 2
[10:03:00] 2 3
[10:03:30] 4 1
[10:04:00] 2 11 60
[10:04:00] 4 2
[10:05:00] 1 3
[10:05:00] 2 11 50
[10:05:00] 4 4
[10:06:00] 1 4
[10:06:00] 4 9 network issue
[10:07:00] 1 3
[10:08:00] 1 11 20
[10:10:00] 1 3
[10:12:00] 1 4
[10:12:00] 1 6
[10:14:00] 1 10 50
[10:20:00] 1 7
[10:25:00] 1 8
[10:40:00] 5 1
[10:45:00] 5 2
[11:00:00] 5 3`)

	cfg, err := parser.ParseConfig(configPath)
	require.NoError(t, err)

	eventsFile, err := os.Open(eventsPath)
	require.NoError(t, err)
	defer eventsFile.Close()

	events, err := parser.ParseEvents(eventsFile)
	require.NoError(t, err)

	processor := engine.New(cfg)
	for _, event := range events {
		processor.Handle(event)
	}

	require.Equal(t, strings.Join([]string{
		"[09:59:00] Player [1] registered",
		"[10:00:00] Player [1] entered the dungeon",
		"[10:01:00] Player [2] registered",
		"[10:02:00] Player [1] killed the monster",
		"[10:02:00] Player [2] entered the dungeon",
		"[10:02:30] Player [3] is disqualified",
		"[10:03:00] Player [2] killed the monster",
		"[10:03:30] Player [4] registered",
		"[10:04:00] Player [2] recieved [60] of damage",
		"[10:04:00] Player [4] entered the dungeon",
		"[10:05:00] Player [1] killed the monster",
		"[10:05:00] Player [2] recieved [50] of damage",
		"[10:05:00] Player [2] is dead",
		"[10:05:00] Player [4] makes imposible move [4]",
		"[10:06:00] Player [1] went to the next floor",
		"[10:06:00] Player [4] cannot continue due to [network issue]",
		"[10:07:00] Player [1] killed the monster",
		"[10:08:00] Player [1] recieved [20] of damage",
		"[10:10:00] Player [1] killed the monster",
		"[10:12:00] Player [1] went to the next floor",
		"[10:12:00] Player [1] entered the boss's floor",
		"[10:14:00] Player [1] has restored [50] of health",
		"[10:20:00] Player [1] killed the boss",
		"[10:25:00] Player [1] left the dungeon",
		"[10:40:00] Player [5] registered",
		"[10:45:00] Player [5] entered the dungeon",
	}, "\n"), formatLogs(processor))

	require.Equal(t, `Final report:
[SUCCESS] 1 [00:25:00, 00:04:30, 00:08:00] HP:100
[FAIL] 2 [00:03:00, 00:00:00, 00:00:00] HP:0
[DISQUAL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100
[DISQUAL] 4 [00:02:00, 00:00:00, 00:00:00] HP:100
[FAIL] 5 [00:15:00, 00:00:00, 00:00:00] HP:100`, output.FormatReport(processor.Report()))
}

func writeTempFile(t *testing.T, name string, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	err := os.WriteFile(path, []byte(content), 0o600)
	require.NoError(t, err)

	return path
}

func formatLogs(processor *engine.Engine) string {
	logs := processor.Logs()
	lines := make([]string, 0, len(logs))

	for _, log := range logs {
		lines = append(lines, output.FormatEvent(log))
	}

	return strings.Join(lines, "\n")
}
