package domain

import "errors"

var (
	ErrInvalidConfig      = errors.New("invalid dungeon config")
	ErrInvalidEventFormat = errors.New("invalid event format")
	ErrUnknownEventType   = errors.New("unknown event type")
	ErrImpossibleMove     = errors.New("impossible move")
	ErrTerminalPlayer     = errors.New("player is in terminal state")
)
