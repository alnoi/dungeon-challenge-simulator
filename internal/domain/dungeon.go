package domain

import "time"

type DungeonConfig struct {
	Floors   int
	Monsters int
	OpenAt   time.Duration
	Duration time.Duration
}

func (c DungeonConfig) CloseAt() time.Duration {
	return c.OpenAt + c.Duration
}
