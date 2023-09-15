package bus

import (
	log "mylife-home-common/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type Client struct {
}

func init() {
	logger := log.CreateLogger("paho:mqtt")

	mqtt.ERROR = &mqttLogger{logger, logrus.ErrorLevel}
	mqtt.CRITICAL = &mqttLogger{logger, logrus.ErrorLevel}
	mqtt.DEBUG = &mqttLogger{logger, logrus.DebugLevel}
	mqtt.WARN = &mqttLogger{logger, logrus.WarnLevel}
}

type mqttLogger struct {
	logger *logrus.Entry
	level  logrus.Level
}

func (ml *mqttLogger) Println(v ...interface{}) {
	ml.logger.Logln(ml.level, v...)
}

func (ml *mqttLogger) Printf(format string, v ...interface{}) {
	ml.logger.Logf(ml.level, format, v...)
}
