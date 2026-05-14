package output

import (
	"fmt"

	"dungeon-event-processor/internal/domain"
	"dungeon-event-processor/internal/timeutil"
)

func FormatEvent(event domain.OutputEvent) string {
	timestamp := timeutil.FormatTimestamp(event.Timestamp)

	switch event.OutgoingType {
	case domain.OutputEventDisqualified:
		return fmt.Sprintf("%s Player [%d] is disqualified", timestamp, event.PlayerID)
	case domain.OutputEventDead:
		return fmt.Sprintf("%s Player [%d] is dead", timestamp, event.PlayerID)
	case domain.OutputEventImpossibleMove:
		return fmt.Sprintf("%s Player [%d] makes imposible move [%d]", timestamp, event.PlayerID, event.RelatedEventType)
	}

	switch event.EventType {
	case domain.EventTypeRegister:
		return fmt.Sprintf("%s Player [%d] registered", timestamp, event.PlayerID)
	case domain.EventTypeEnterDungeon:
		return fmt.Sprintf("%s Player [%d] entered the dungeon", timestamp, event.PlayerID)
	case domain.EventTypeKillMonster:
		return fmt.Sprintf("%s Player [%d] killed the monster", timestamp, event.PlayerID)
	case domain.EventTypeNextFloor:
		return fmt.Sprintf("%s Player [%d] went to the next floor", timestamp, event.PlayerID)
	case domain.EventTypePreviousFloor:
		return fmt.Sprintf("%s Player [%d] went to the previous floor", timestamp, event.PlayerID)
	case domain.EventTypeEnterBossFloor:
		return fmt.Sprintf("%s Player [%d] entered the boss's floor", timestamp, event.PlayerID)
	case domain.EventTypeKillBoss:
		return fmt.Sprintf("%s Player [%d] killed the boss", timestamp, event.PlayerID)
	case domain.EventTypeLeaveDungeon:
		return fmt.Sprintf("%s Player [%d] left the dungeon", timestamp, event.PlayerID)
	case domain.EventTypeCannotContinue:
		return fmt.Sprintf("%s Player [%d] cannot continue due to [%s]", timestamp, event.PlayerID, event.Extra)
	case domain.EventTypeHeal:
		return fmt.Sprintf("%s Player [%d] has restored [%s] of health", timestamp, event.PlayerID, event.Extra)
	case domain.EventTypeDamage:
		return fmt.Sprintf("%s Player [%d] recieved [%s] of damage", timestamp, event.PlayerID, event.Extra)
	default:
		return fmt.Sprintf("%s Player [%d] produced unknown event", timestamp, event.PlayerID)
	}
}
