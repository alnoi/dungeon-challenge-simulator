package engine_test

import (
	"testing"
	"time"

	"dungeon-event-processor/internal/domain"
	"dungeon-event-processor/internal/engine"
	"dungeon-event-processor/internal/output"

	"github.com/stretchr/testify/require"
)

func TestEngineDisqualifiesUnregisteredPlayer(t *testing.T) {
	processor := newTestEngine()

	processor.Handle(event(14*time.Hour+10*time.Minute, 3, domain.EventTypeEnterDungeon, ""))

	require.Equal(t, []string{
		"[14:10:00] Player [3] is disqualified",
	}, formattedLogs(processor))

	report := processor.Report()
	require.Len(t, report, 1)
	require.Equal(t, domain.FinalStateDisqual, report[0].State)
	require.Equal(t, 3, report[0].PlayerID)
	require.Equal(t, domain.MaxHealth, report[0].Health)
}

func TestEngineRejectsMoveNextBeforeFloorIsCleared(t *testing.T) {
	processor := newTestEngine()

	processor.Handle(event(14*time.Hour, 1, domain.EventTypeRegister, ""))
	processor.Handle(event(14*time.Hour+time.Minute, 1, domain.EventTypeEnterDungeon, ""))
	processor.Handle(event(14*time.Hour+2*time.Minute, 1, domain.EventTypeNextFloor, ""))

	require.Equal(t, []string{
		"[14:00:00] Player [1] registered",
		"[14:01:00] Player [1] entered the dungeon",
		"[14:02:00] Player [1] makes imposible move [4]",
	}, formattedLogs(processor))
}

func TestEngineCapsHealAtMaxHealth(t *testing.T) {
	processor := newTestEngine()

	processor.Handle(event(14*time.Hour, 1, domain.EventTypeRegister, ""))
	processor.Handle(event(14*time.Hour+time.Minute, 1, domain.EventTypeEnterDungeon, ""))
	processor.Handle(event(14*time.Hour+2*time.Minute, 1, domain.EventTypeDamage, "30"))
	processor.Handle(event(14*time.Hour+3*time.Minute, 1, domain.EventTypeHeal, "80"))

	report := processor.Report()
	require.Len(t, report, 1)
	require.Equal(t, domain.MaxHealth, report[0].Health)
}

func TestEngineEmitsDeathWhenDamageDropsHealthToZero(t *testing.T) {
	processor := newTestEngine()

	processor.Handle(event(14*time.Hour, 1, domain.EventTypeRegister, ""))
	processor.Handle(event(14*time.Hour+time.Minute, 1, domain.EventTypeEnterDungeon, ""))
	processor.Handle(event(14*time.Hour+2*time.Minute, 1, domain.EventTypeDamage, "120"))

	require.Equal(t, []string{
		"[14:00:00] Player [1] registered",
		"[14:01:00] Player [1] entered the dungeon",
		"[14:02:00] Player [1] recieved [120] of damage",
		"[14:02:00] Player [1] is dead",
	}, formattedLogs(processor))

	report := processor.Report()
	require.Len(t, report, 1)
	require.Equal(t, domain.FinalStateFail, report[0].State)
	require.Equal(t, 0, report[0].Health)
}

func TestEngineRejectsBossFloorBeforeRegularFloorsAreCleared(t *testing.T) {
	processor := newTestEngine()

	processor.Handle(event(14*time.Hour, 1, domain.EventTypeRegister, ""))
	processor.Handle(event(14*time.Hour+time.Minute, 1, domain.EventTypeEnterDungeon, ""))
	processor.Handle(event(14*time.Hour+2*time.Minute, 1, domain.EventTypeEnterBossFloor, ""))

	require.Equal(t, []string{
		"[14:00:00] Player [1] registered",
		"[14:01:00] Player [1] entered the dungeon",
		"[14:02:00] Player [1] makes imposible move [6]",
	}, formattedLogs(processor))
}

func TestEngineExpiresActivePlayerWhenDungeonCloses(t *testing.T) {
	processor := newTestEngine()

	processor.Handle(event(14*time.Hour, 1, domain.EventTypeRegister, ""))
	processor.Handle(event(14*time.Hour, 1, domain.EventTypeEnterDungeon, ""))
	processor.Handle(event(15*time.Hour, 1, domain.EventTypeDamage, "10"))

	require.Equal(t, []string{
		"[14:00:00] Player [1] registered",
		"[14:00:00] Player [1] entered the dungeon",
	}, formattedLogs(processor))

	report := processor.Report()
	require.Len(t, report, 1)
	require.Equal(t, domain.FinalStateFail, report[0].State)
	require.Equal(t, time.Hour, report[0].TotalTime)
	require.Equal(t, domain.MaxHealth, report[0].Health)
}

func newTestEngine() *engine.Engine {
	return engine.New(domain.DungeonConfig{
		Floors:   2,
		Monsters: 2,
		OpenAt:   14 * time.Hour,
		Duration: time.Hour,
	})
}

func event(timestamp time.Duration, playerID int, eventType domain.EventType, extra string) domain.Event {
	return domain.Event{
		Timestamp: timestamp,
		PlayerID:  playerID,
		Type:      eventType,
		Extra:     extra,
	}
}

func formattedLogs(processor *engine.Engine) []string {
	logs := processor.Logs()
	formatted := make([]string, 0, len(logs))

	for _, log := range logs {
		formatted = append(formatted, output.FormatEvent(log))
	}

	return formatted
}
