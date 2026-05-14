package domain

import "time"

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeRegister
	EventTypeEnterDungeon
	EventTypeKillMonster
	EventTypeNextFloor
	EventTypePreviousFloor
	EventTypeEnterBossFloor
	EventTypeKillBoss
	EventTypeLeaveDungeon
	EventTypeCannotContinue
	EventTypeHeal
	EventTypeDamage
)

func (t EventType) IsKnown() bool {
	return t >= EventTypeRegister && t <= EventTypeDamage
}

type OutputEventType int

const (
	OutputEventTypeUnknown OutputEventType = 0
)

const (
	OutputEventDisqualified OutputEventType = 31 + iota
	OutputEventDead
	OutputEventImpossibleMove
)

type Event struct {
	Timestamp time.Duration
	PlayerID  int
	Type      EventType
	Extra     string
}

type OutputEvent struct {
	Timestamp        time.Duration
	PlayerID         int
	EventType        EventType
	OutgoingType     OutputEventType
	Extra            string
	RelatedEventType EventType
}
