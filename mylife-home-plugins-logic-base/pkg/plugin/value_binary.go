//go:generate cd ../../../mylife-home-core-generator && go run cmd/main.go -- $GOFILE
package plugin

// @Plugin(description="binary value description" usage=logic)
type ValueBinary struct {
	// @State
	Value State[bool]
}

// @Action
func (this *ValueBinary) SetValue(arg bool) {
	this.Value.Set(arg)
}
