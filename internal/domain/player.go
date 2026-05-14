package domain

import "time"

const MaxHealth = 100

type PlayerState int

const (
	PlayerStateUnknown PlayerState = iota
	PlayerStateRegistered
	PlayerStateInDungeon
	PlayerStateOnBossFloor
	PlayerStateCompleted
	PlayerStateLeft
	PlayerStateDead
	PlayerStateDisqualified
	PlayerStateExpired
)

type Player struct {
	ID           int
	State        PlayerState
	Health       int
	CurrentFloor int
	Entered      bool

	EnteredAt  time.Duration
	FinishedAt time.Duration

	FloorEnteredAt       time.Duration
	FloorTimerRunning    bool
	FloorMonsterKills    map[int]int
	ClearedFloors        map[int]bool
	FloorActiveDurations map[int]time.Duration
	FloorClearDurations  []time.Duration

	BossStartedAt    time.Duration
	BossTimerRunning bool
	BossDuration     time.Duration
	BossKilled       bool
}

func NewPlayer(id int) *Player {
	return &Player{
		ID:                   id,
		State:                PlayerStateUnknown,
		Health:               MaxHealth,
		FloorMonsterKills:    make(map[int]int),
		ClearedFloors:        make(map[int]bool),
		FloorActiveDurations: make(map[int]time.Duration),
	}
}
