package engine

import (
	"sort"
	"time"

	"dungeon-event-processor/internal/domain"
)

type Engine struct {
	config  domain.DungeonConfig
	players map[int]*domain.Player
	logs    []domain.OutputEvent
}

func New(cfg domain.DungeonConfig) *Engine {
	return &Engine{
		config:  cfg,
		players: make(map[int]*domain.Player),
		logs:    make([]domain.OutputEvent, 0),
	}
}

func (e *Engine) Handle(event domain.Event) {
	player := e.ensurePlayer(event.PlayerID)

	if e.isTerminal(player) {
		return
	}

	if e.hasDungeonClosedFor(player, event.Timestamp) {
		e.expirePlayer(player)
		return
	}

	if player.State == domain.PlayerStateUnknown && event.Type != domain.EventTypeRegister {
		player.State = domain.PlayerStateDisqualified
		e.appendOutgoing(event.Timestamp, player.ID, domain.OutputEventDisqualified, event.Type)
		return
	}

	switch event.Type {
	case domain.EventTypeRegister:
		e.handleRegister(player, event)
	case domain.EventTypeEnterDungeon:
		e.handleEnterDungeon(player, event)
	case domain.EventTypeKillMonster:
		e.handleKillMonster(player, event)
	case domain.EventTypeNextFloor:
		e.handleNextFloor(player, event)
	case domain.EventTypePreviousFloor:
		e.handlePreviousFloor(player, event)
	case domain.EventTypeEnterBossFloor:
		e.handleEnterBossFloor(player, event)
	case domain.EventTypeKillBoss:
		e.handleKillBoss(player, event)
	case domain.EventTypeLeaveDungeon:
		e.handleLeaveDungeon(player, event)
	case domain.EventTypeCannotContinue:
		e.handleCannotContinue(player, event)
	case domain.EventTypeHeal:
		e.handleHeal(player, event)
	case domain.EventTypeDamage:
		e.handleDamage(player, event)
	}
}

func (e *Engine) Logs() []domain.OutputEvent {
	logs := make([]domain.OutputEvent, len(e.logs))
	copy(logs, e.logs)
	return logs
}

func (e *Engine) Report() []domain.ReportLine {
	ids := make([]int, 0, len(e.players))
	for id := range e.players {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	report := make([]domain.ReportLine, 0, len(ids))
	for _, id := range ids {
		player := e.players[id]
		finishedAt := e.finishedAt(player)

		report = append(report, domain.ReportLine{
			State:                 e.finalState(player),
			PlayerID:              player.ID,
			TotalTime:             e.totalTime(player, finishedAt),
			AverageFloorClearTime: e.averageFloorClearTime(player),
			BossKillTime:          player.BossDuration,
			Health:                player.Health,
		})
	}

	return report
}

func (e *Engine) ensurePlayer(id int) *domain.Player {
	if player, ok := e.players[id]; ok {
		return player
	}

	player := domain.NewPlayer(id)
	e.players[id] = player
	return player
}

func (e *Engine) appendIncoming(event domain.Event) {
	e.logs = append(e.logs, domain.OutputEvent{
		Timestamp: event.Timestamp,
		PlayerID:  event.PlayerID,
		EventType: event.Type,
		Extra:     event.Extra,
	})
}

func (e *Engine) appendOutgoing(timestamp time.Duration, playerID int, eventType domain.OutputEventType, related domain.EventType) {
	e.logs = append(e.logs, domain.OutputEvent{
		Timestamp:        timestamp,
		PlayerID:         playerID,
		OutgoingType:     eventType,
		RelatedEventType: related,
	})
}

func (e *Engine) finalState(player *domain.Player) domain.FinalState {
	if player.State == domain.PlayerStateDisqualified {
		return domain.FinalStateDisqual
	}

	if player.State == domain.PlayerStateDead {
		return domain.FinalStateFail
	}

	if player.BossKilled && e.allRegularFloorsCleared(player) {
		return domain.FinalStateSuccess
	}

	return domain.FinalStateFail
}

func (e *Engine) regularFloors() int {
	if e.config.Floors <= 0 {
		return 0
	}

	return e.config.Floors - 1
}

func (e *Engine) hasDungeonClosedFor(player *domain.Player, timestamp time.Duration) bool {
	return e.isInChallenge(player) && timestamp >= e.config.CloseAt()
}

func (e *Engine) expirePlayer(player *domain.Player) {
	closeAt := e.config.CloseAt()
	e.stopFloorTimer(player, closeAt)
	player.State = domain.PlayerStateExpired
	player.FinishedAt = closeAt
}

func (e *Engine) isInChallenge(player *domain.Player) bool {
	switch player.State {
	case domain.PlayerStateInDungeon, domain.PlayerStateOnBossFloor, domain.PlayerStateCompleted:
		return true
	default:
		return false
	}
}

func (e *Engine) finishedAt(player *domain.Player) time.Duration {
	if player.FinishedAt > 0 {
		return player.FinishedAt
	}

	if e.isInChallenge(player) {
		return e.config.CloseAt()
	}

	return 0
}

func (e *Engine) totalTime(player *domain.Player, finishedAt time.Duration) time.Duration {
	if !player.Entered || finishedAt <= player.EnteredAt {
		return 0
	}

	return finishedAt - player.EnteredAt
}

func (e *Engine) averageFloorClearTime(player *domain.Player) time.Duration {
	if len(player.FloorClearDurations) == 0 {
		return 0
	}

	var total time.Duration
	for _, duration := range player.FloorClearDurations {
		total += duration
	}

	return total / time.Duration(len(player.FloorClearDurations))
}

func (e *Engine) allRegularFloorsCleared(player *domain.Player) bool {
	for floor := 1; floor <= e.regularFloors(); floor++ {
		if !e.isRegularFloorCleared(player, floor) {
			return false
		}
	}

	return true
}

func (e *Engine) isRegularFloorCleared(player *domain.Player, floor int) bool {
	if floor < 1 || floor > e.regularFloors() {
		return false
	}

	if e.config.Monsters == 0 {
		return true
	}

	return player.ClearedFloors[floor]
}

func (e *Engine) startFloorTimer(player *domain.Player, timestamp time.Duration) {
	if player.CurrentFloor < 1 || player.CurrentFloor > e.regularFloors() {
		player.FloorTimerRunning = false
		return
	}

	if e.isRegularFloorCleared(player, player.CurrentFloor) {
		player.FloorTimerRunning = false
		return
	}

	player.FloorEnteredAt = timestamp
	player.FloorTimerRunning = true
}

func (e *Engine) stopFloorTimer(player *domain.Player, timestamp time.Duration) {
	if !player.FloorTimerRunning {
		return
	}

	if timestamp > player.FloorEnteredAt {
		player.FloorActiveDurations[player.CurrentFloor] += timestamp - player.FloorEnteredAt
	}

	player.FloorTimerRunning = false
}

func (e *Engine) clearCurrentFloor(player *domain.Player, timestamp time.Duration) {
	e.stopFloorTimer(player, timestamp)
	player.ClearedFloors[player.CurrentFloor] = true
	player.FloorClearDurations = append(player.FloorClearDurations, player.FloorActiveDurations[player.CurrentFloor])
}
