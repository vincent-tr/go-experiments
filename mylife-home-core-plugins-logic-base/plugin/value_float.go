package plugin

import (
	"mylife-home-core-common/definitions"
)

// @Plugin(description="binary value description" usage="logic")
type ValueFloat struct {
	// @Config(description="initial value")
	InitialValue float64

	// @State(description="current value")
	Value definitions.State[float64]
}

// @Action(description="set current value")
func (this *ValueFloat) SetValue(arg float64) {
	this.Value.Set(arg)
}

func (this *ValueFloat) Init() error {
	this.Value.Set(this.InitialValue)
}

func (this *ValueFloat) Terminate() {
	// Noop
}
