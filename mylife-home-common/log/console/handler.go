// From https://github.com/apex/log/blob/master/handlers/text/text.go
package console

import (
	"fmt"
	"io"
	"sync"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Colors mapping.
var Colors = [...]int{
	log.DebugLevel: gray,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

// Handler implementation.
type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer: w,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]
	loggerName := e.Fields.Get("logger-name").(string)

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintf(h.Writer, "%s - %-25s - \033[%dm%5s\033[0m %s", e.Timestamp.Format("2006-01-02 15:04:05"), loggerName, color, level, e.Message)

	err := e.Fields.Get("error")
	if err != nil {

		fmt.Fprintf(h.Writer, ": Error: %s", err.(error).Error())

		if st, ok := err.(stackTracer); ok {
			fmt.Fprintf(h.Writer, "%+v", st.StackTrace())
		}
	}

	fmt.Fprintln(h.Writer)

	return nil
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
