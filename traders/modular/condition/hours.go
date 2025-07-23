package condition

import (
	"fmt"
	"go-experiments/common"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

func Hours(startHour, endHour int) Condition {
	return &hoursCondition{
		startHour: startHour,
		endHour:   endHour,
	}
}

type hoursCondition struct {
	startHour int
	endHour   int
}

func (h *hoursCondition) Execute(ctx context.TraderContext) bool {
	hour := ctx.Timestamp().Hour()
	return hour >= h.startHour && hour < h.endHour
}

func (h *hoursCondition) Format() *formatter.FormatterNode {
	return formatter.Format("Hours",
		formatter.Format(fmt.Sprintf("StartHour: %d", h.startHour)),
		formatter.Format(fmt.Sprintf("EndHour: %d", h.endHour)),
	)
}

func Session(session *common.Session) Condition {
	return &sessionCondition{
		session: session,
	}
}

type sessionCondition struct {
	session *common.Session
}

func (s *sessionCondition) Execute(ctx context.TraderContext) bool {
	return s.session.IsOpen(ctx.Timestamp())
}

func (s *sessionCondition) Format() *formatter.FormatterNode {
	return formatter.Format("Session",
		formatter.Format(fmt.Sprintf("Session: %s", s.session.String())),
	)
}
