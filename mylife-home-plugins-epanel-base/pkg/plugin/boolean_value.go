package plugin

type BooleanValue struct {
	InitialValue Config[bool]

	Value State[bool]

	Toggle Action[bool]
	On Action[bool]
	Off Action[bool]
}

func (this *BooleanValue) Init() error {

}

func (this *BooleanValue) toggle(arg bool) {
	this.Value.Set(!this.Value.Get())
}

func (this *BooleanValue) on(arg bool) {
	this.Value.Set(true)
}

func (this *BooleanValue) off(arg bool) {
	this.Value.Set(false)
}

func init() {
	registry.RegisterPlugin[BooleanValue]();
}