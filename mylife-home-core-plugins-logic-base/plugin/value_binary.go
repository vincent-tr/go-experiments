package plugin

import (
	"mylife-home-core-common/definitions"
)

// @Plugin(description="binary value description" usage="logic")
type ValueBinary struct {
	// @Config
	InitialValue bool

	// @State
	Value definitions.State[bool]
}

// @Action
func (this *ValueBinary) SetValue(arg bool) {
	this.Value.Set(arg)
}

func (this *ValueBinary) Init() error {
	this.Value.Set(this.InitialValue)
}

func (this *ValueBinary) Terminate() {
	// Noop
}
