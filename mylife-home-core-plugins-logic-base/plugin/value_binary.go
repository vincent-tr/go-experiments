package plugin

import (
	"mylife-home-core-common/definitions"
)

// @Plugin(description="binary value description" usage="logic")
type ValueBinary struct {
	// @State
	Value definitions.State[bool]
}

// @Action
func (this *ValueBinary) SetValue(arg bool) {
	this.Value.Set(arg)
}
