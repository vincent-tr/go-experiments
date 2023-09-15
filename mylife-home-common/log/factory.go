package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Logger = *log.Entry
type Fields = log.Fields

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}

func CreateLogger(name string) Logger {
	return log.WithFields(log.Fields{
		"logger-name": name,
	})
}
