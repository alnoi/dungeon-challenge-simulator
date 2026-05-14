package output

import (
	"testing"
	"time"

	"dungeon-event-processor/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestFormatReport(t *testing.T) {
	lines := []domain.ReportLine{
		{
			State:                 domain.FinalStateSuccess,
			PlayerID:              1,
			TotalTime:             24 * time.Minute,
			AverageFloorClearTime: 5 * time.Minute,
			BossKillTime:          11 * time.Minute,
			Health:                35,
		},
		{
			State:    domain.FinalStateDisqual,
			PlayerID: 3,
			Health:   100,
		},
	}

	require.Equal(t, `Final report:
[SUCCESS] 1 [00:24:00, 00:05:00, 00:11:00] HP:35
[DISQUAL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100`, FormatReport(lines))
}
