package engine

import "dungeon-event-processor/internal/domain"

func (e *Engine) isTerminal(player *domain.Player) bool {
	switch player.State {
	case domain.PlayerStateLeft, domain.PlayerStateDead, domain.PlayerStateDisqualified, domain.PlayerStateExpired:
		return true
	default:
		return false
	}
}

func (e *Engine) canEnterDungeon(player *domain.Player, event domain.Event) bool {
	if player.State != domain.PlayerStateRegistered {
		return false
	}

	return event.Timestamp >= e.config.OpenAt && event.Timestamp < e.config.CloseAt()
}

func (e *Engine) canKillMonster(player *domain.Player, event domain.Event) bool {
	if player.State != domain.PlayerStateInDungeon {
		return false
	}

	if player.CurrentFloor < 1 || player.CurrentFloor > e.regularFloors() {
		return false
	}

	if e.isRegularFloorCleared(player, player.CurrentFloor) {
		return false
	}

	return player.FloorMonsterKills[player.CurrentFloor] < e.config.Monsters
}

func (e *Engine) canMoveNext(player *domain.Player, event domain.Event) bool {
	if player.State != domain.PlayerStateInDungeon {
		return false
	}

	if player.CurrentFloor < 1 || player.CurrentFloor >= e.config.Floors {
		return false
	}

	if player.CurrentFloor <= e.regularFloors() && !e.isRegularFloorCleared(player, player.CurrentFloor) {
		return false
	}

	return true
}

func (e *Engine) canMovePrevious(player *domain.Player, event domain.Event) bool {
	if player.State != domain.PlayerStateInDungeon && player.State != domain.PlayerStateOnBossFloor {
		return false
	}

	return player.CurrentFloor > 1
}

func (e *Engine) canEnterBossFloor(player *domain.Player, event domain.Event) bool {
	if player.State != domain.PlayerStateInDungeon {
		return false
	}

	if player.CurrentFloor != e.config.Floors {
		return false
	}

	if player.BossKilled {
		return false
	}

	return e.allRegularFloorsCleared(player)
}

func (e *Engine) canKillBoss(player *domain.Player, event domain.Event) bool {
	return player.State == domain.PlayerStateOnBossFloor &&
		player.CurrentFloor == e.config.Floors &&
		!player.BossKilled
}
