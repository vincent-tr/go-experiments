package conditions

import (
	"encoding/json"
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

func init() {
	jsonParsers.RegisterParser("hours", func(arg json.RawMessage) (Condition, error) {
		var hours struct {
			StartHour int `json:"startHour"`
			EndHour   int `json:"endHour"`
		}
		if err := json.Unmarshal(arg, &hours); err != nil {
			return nil, fmt.Errorf("failed to parse hours condition: %w", err)
		}

		return Hours(hours.StartHour, hours.EndHour), nil
	})
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

func init() {
	jsonParsers.RegisterParser("session", func(arg json.RawMessage) (Condition, error) {
		var sessionName string
		if err := json.Unmarshal(arg, &sessionName); err != nil {
			return nil, fmt.Errorf("failed to parse session condition: %w", err)
		}

		var session *common.Session
		// TODO: more dynamic
		switch sessionName {
		case "London":
			session = common.LondonSession
		case "New York":
			session = common.NYSession
		default:
			return nil, fmt.Errorf("unknown session: %s", sessionName)
		}

		return Session(session), nil
	})
}
