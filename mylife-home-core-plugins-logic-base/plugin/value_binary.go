package plugin

import (
	"mylife-home-core-common/definitions"
)

// @Plugin(description="binary value description" usage="logic")
type ValueBinary struct {
	// @Config(description="initial value")
	InitialValue bool

	// @State(description="current value")
	Value definitions.State[bool]
}

// @Action(description="set current value")
func (this *ValueBinary) SetValue(arg bool) {
	this.Value.Set(arg)
}

func (this *ValueBinary) Init() error {
	this.Value.Set(this.InitialValue)
}

func (this *ValueBinary) Terminate() {
	// Noop
}
