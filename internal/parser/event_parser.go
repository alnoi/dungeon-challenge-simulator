package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"dungeon-event-processor/internal/domain"
	"dungeon-event-processor/internal/timeutil"
)

func ParseEvent(line string) (domain.Event, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return domain.Event{}, fmt.Errorf("%w: %q", domain.ErrInvalidEventFormat, line)
	}

	timestamp, err := timeutil.ParseClock(parts[0])
	if err != nil {
		return domain.Event{}, fmt.Errorf("%w: invalid timestamp: %v", domain.ErrInvalidEventFormat, err)
	}

	playerID, err := strconv.Atoi(parts[1])
	if err != nil {
		return domain.Event{}, fmt.Errorf("%w: invalid player id: %v", domain.ErrInvalidEventFormat, err)
	}

	eventID, err := strconv.Atoi(parts[2])
	if err != nil {
		return domain.Event{}, fmt.Errorf("%w: invalid event id: %v", domain.ErrInvalidEventFormat, err)
	}

	eventType := domain.EventType(eventID)
	if !eventType.IsKnown() {
		return domain.Event{}, fmt.Errorf("%w: %d", domain.ErrUnknownEventType, eventID)
	}

	return domain.Event{
		Timestamp: timestamp,
		PlayerID:  playerID,
		Type:      eventType,
		Extra:     strings.Join(parts[3:], " "),
	}, nil
}

func ParseEvents(reader io.Reader) ([]domain.Event, error) {
	scanner := bufio.NewScanner(reader)
	events := make([]domain.Event, 0)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event, err := ParseEvent(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}

	return events, nil
}
