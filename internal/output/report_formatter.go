package output

import (
	"fmt"
	"strings"

	"dungeon-event-processor/internal/domain"
	"dungeon-event-processor/internal/timeutil"
)

func FormatReport(lines []domain.ReportLine) string {
	var builder strings.Builder
	builder.WriteString("Final report:")

	for _, line := range lines {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf(
			"[%s] %d [%s, %s, %s] HP:%d",
			line.State.String(),
			line.PlayerID,
			timeutil.FormatDuration(line.TotalTime),
			timeutil.FormatDuration(line.AverageFloorClearTime),
			timeutil.FormatDuration(line.BossKillTime),
			line.Health,
		))
	}

	return builder.String()
}
