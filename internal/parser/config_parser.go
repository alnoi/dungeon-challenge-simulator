package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"dungeon-event-processor/internal/domain"
	"dungeon-event-processor/internal/timeutil"
)

type rawDungeonConfig struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

func ParseConfig(path string) (domain.DungeonConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return domain.DungeonConfig{}, fmt.Errorf("open config: %w", err)
	}
	defer file.Close()

	var raw rawDungeonConfig
	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		return domain.DungeonConfig{}, fmt.Errorf("%w: %v", domain.ErrInvalidConfig, err)
	}

	openAt, err := timeutil.ParseClock(raw.OpenAt)
	if err != nil {
		return domain.DungeonConfig{}, fmt.Errorf("%w: invalid OpenAt: %v", domain.ErrInvalidConfig, err)
	}

	if raw.Floors < 1 {
		return domain.DungeonConfig{}, fmt.Errorf("%w: Floors must be positive", domain.ErrInvalidConfig)
	}

	if raw.Monsters < 0 {
		return domain.DungeonConfig{}, fmt.Errorf("%w: Monsters cannot be negative", domain.ErrInvalidConfig)
	}

	if raw.Duration <= 0 {
		return domain.DungeonConfig{}, fmt.Errorf("%w: Duration must be positive", domain.ErrInvalidConfig)
	}

	return domain.DungeonConfig{
		Floors:   raw.Floors,
		Monsters: raw.Monsters,
		OpenAt:   openAt,
		Duration: time.Duration(raw.Duration) * time.Hour,
	}, nil
}
