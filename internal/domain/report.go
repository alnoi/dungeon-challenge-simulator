package domain

import "time"

type FinalState int

const (
	FinalStateUnknown FinalState = iota
	FinalStateSuccess
	FinalStateFail
	FinalStateDisqual
)

func (s FinalState) String() string {
	switch s {
	case FinalStateSuccess:
		return "SUCCESS"
	case FinalStateFail:
		return "FAIL"
	case FinalStateDisqual:
		return "DISQUAL"
	default:
		return "UNKNOWN"
	}
}

type ReportLine struct {
	State                 FinalState
	PlayerID              int
	TotalTime             time.Duration
	AverageFloorClearTime time.Duration
	BossKillTime          time.Duration
	Health                int
}
