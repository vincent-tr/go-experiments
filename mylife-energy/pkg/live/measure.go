package live

import (
	"mylife-tools-server/services/io/serialization"
	"time"
)

type Measure struct {
	id        string
	sensor    string
	timestamp time.Time
	value     float64
}

func (measure *Measure) Id() string {
	return measure.id
}

func (measure *Measure) Sensor() string {
	return measure.sensor
}

func (measure *Measure) Timestamp() time.Time {
	return measure.timestamp
}

func (measure *Measure) Value() float64 {
	return measure.value
}

func (measure *Measure) Marshal() (interface{}, error) {
	helper := serialization.NewStructMarshallerHelper()

	helper.Add("_id", measure.id)
	helper.Add("sensor", measure.sensor)
	helper.Add("timestamp", measure.timestamp)
	helper.Add("value", measure.value)

	return helper.Build()
}
