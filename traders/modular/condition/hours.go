package condition

import (
	"fmt"
	"go-experiments/common"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

func Hours(startHour, endHour int) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			hour := ctx.Timestamp().Hour()
			return hour >= startHour && hour < endHour
		},
		func() *formatter.FormatterNode {
			return formatter.Format("Hours",
				formatter.Format(fmt.Sprintf("StartHour: %d", startHour)),
				formatter.Format(fmt.Sprintf("EndHour: %d", endHour)),
			)
		},
	)
}

func Session(session *common.Session) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return session.IsOpen(ctx.Timestamp())
		},
		func() *formatter.FormatterNode {
			return formatter.Format(fmt.Sprintf("Session: %s", session.String()))
		},
	)
}
