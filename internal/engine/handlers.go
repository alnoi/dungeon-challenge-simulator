package engine

import (
	"strconv"
	"strings"

	"dungeon-event-processor/internal/domain"
)

func (e *Engine) handleRegister(player *domain.Player, event domain.Event) {
	if player.State != domain.PlayerStateUnknown {
		e.appendImpossibleMove(event)
		return
	}

	player.State = domain.PlayerStateRegistered
	player.Health = domain.MaxHealth
	e.appendIncoming(event)
}

func (e *Engine) handleEnterDungeon(player *domain.Player, event domain.Event) {
	if !e.canEnterDungeon(player, event) {
		e.appendImpossibleMove(event)
		return
	}

	player.State = domain.PlayerStateInDungeon
	player.CurrentFloor = 1
	player.Entered = true
	player.EnteredAt = event.Timestamp
	e.startFloorTimer(player, event.Timestamp)
	e.appendIncoming(event)
}

func (e *Engine) handleKillMonster(player *domain.Player, event domain.Event) {
	if !e.canKillMonster(player, event) {
		e.appendImpossibleMove(event)
		return
	}

	player.FloorMonsterKills[player.CurrentFloor]++
	if player.FloorMonsterKills[player.CurrentFloor] >= e.config.Monsters {
		e.clearCurrentFloor(player, event.Timestamp)
	}

	e.appendIncoming(event)
}

func (e *Engine) handleNextFloor(player *domain.Player, event domain.Event) {
	if !e.canMoveNext(player, event) {
		e.appendImpossibleMove(event)
		return
	}

	e.stopFloorTimer(player, event.Timestamp)
	player.CurrentFloor++
	e.startFloorTimer(player, event.Timestamp)
	e.appendIncoming(event)
}

func (e *Engine) handlePreviousFloor(player *domain.Player, event domain.Event) {
	if !e.canMovePrevious(player, event) {
		e.appendImpossibleMove(event)
		return
	}

	e.stopFloorTimer(player, event.Timestamp)
	player.CurrentFloor--
	if player.State == domain.PlayerStateOnBossFloor {
		player.State = domain.PlayerStateInDungeon
	}
	e.startFloorTimer(player, event.Timestamp)
	e.appendIncoming(event)
}

func (e *Engine) handleEnterBossFloor(player *domain.Player, event domain.Event) {
	if !e.canEnterBossFloor(player, event) {
		e.appendImpossibleMove(event)
		return
	}

	player.State = domain.PlayerStateOnBossFloor
	player.BossStartedAt = event.Timestamp
	player.BossTimerRunning = true
	e.appendIncoming(event)
}

func (e *Engine) handleKillBoss(player *domain.Player, event domain.Event) {
	if !e.canKillBoss(player, event) {
		e.appendImpossibleMove(event)
		return
	}

	player.BossKilled = true
	player.State = domain.PlayerStateCompleted
	player.BossTimerRunning = false
	if event.Timestamp > player.BossStartedAt {
		player.BossDuration = event.Timestamp - player.BossStartedAt
	}
	e.appendIncoming(event)
}

func (e *Engine) handleLeaveDungeon(player *domain.Player, event domain.Event) {
	if !e.isInChallenge(player) {
		e.appendImpossibleMove(event)
		return
	}

	e.stopFloorTimer(player, event.Timestamp)
	player.State = domain.PlayerStateLeft
	player.FinishedAt = event.Timestamp
	e.appendIncoming(event)
}

func (e *Engine) handleCannotContinue(player *domain.Player, event domain.Event) {
	if player.State == domain.PlayerStateUnknown {
		player.State = domain.PlayerStateDisqualified
		e.appendOutgoing(event.Timestamp, player.ID, domain.OutputEventDisqualified, event.Type)
		return
	}

	if e.isInChallenge(player) {
		e.stopFloorTimer(player, event.Timestamp)
		player.FinishedAt = event.Timestamp
	}

	player.State = domain.PlayerStateDisqualified
	e.appendIncoming(event)
}

func (e *Engine) handleHeal(player *domain.Player, event domain.Event) {
	value, ok := parsePositiveInt(event.Extra)
	if !ok || !e.isInChallenge(player) {
		e.appendImpossibleMove(event)
		return
	}

	player.Health += value
	if player.Health > domain.MaxHealth {
		player.Health = domain.MaxHealth
	}

	e.appendIncoming(event)
}

func (e *Engine) handleDamage(player *domain.Player, event domain.Event) {
	value, ok := parsePositiveInt(event.Extra)
	if !ok || !e.isInChallenge(player) {
		e.appendImpossibleMove(event)
		return
	}

	player.Health -= value
	if player.Health < 0 {
		player.Health = 0
	}

	e.appendIncoming(event)
	if player.Health == 0 {
		e.stopFloorTimer(player, event.Timestamp)
		player.State = domain.PlayerStateDead
		player.FinishedAt = event.Timestamp
		e.appendOutgoing(event.Timestamp, player.ID, domain.OutputEventDead, event.Type)
	}
}

func parsePositiveInt(value string) (int, bool) {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < 0 {
		return 0, false
	}

	return parsed, true
}

func (e *Engine) appendImpossibleMove(event domain.Event) {
	e.appendOutgoing(event.Timestamp, event.PlayerID, domain.OutputEventImpossibleMove, event.Type)
}
