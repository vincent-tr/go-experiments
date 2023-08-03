package plugin

// @Plugin(description='boolean value description' usage=logic)
type ValueBinary struct {
	// @State
	Value State[bool]
}

// @Action
func (this *ValueBinary) SetValue(arg bool) {
	this.Value.Set(arg)
}
