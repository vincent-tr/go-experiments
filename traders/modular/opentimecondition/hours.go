package opentimecondition

import (
	"fmt"
	"go-experiments/common"
	"go-experiments/traders/modular/formatter"
	"time"
)

func Hours(startHour, endHour int) OpenTimeCondition {
	return &hoursCondition{
		startHour: startHour,
		endHour:   endHour,
	}
}

type hoursCondition struct {
	startHour int
	endHour   int
}

func (h *hoursCondition) Execute(timestamp time.Time) bool {
	return timestamp.Hour() >= h.startHour && timestamp.Hour() < h.endHour
}

func (h *hoursCondition) Format() *formatter.FormatterNode {
	return formatter.Format("Hours",
		formatter.Format(fmt.Sprintf("StartHour: %d", h.startHour)),
		formatter.Format(fmt.Sprintf("EndHour: %d", h.endHour)),
	)
}

func Session(session *common.Session) OpenTimeCondition {
	return &sessionCondition{
		session: session,
	}
}

type sessionCondition struct {
	session *common.Session
}

func (s *sessionCondition) Execute(timestamp time.Time) bool {
	return s.session.IsOpen(timestamp)
}

func (s *sessionCondition) Format() *formatter.FormatterNode {
	return formatter.Format("Session",
		formatter.Format(fmt.Sprintf("Session: %s", s.session.String())),
	)
}
