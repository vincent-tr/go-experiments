package impl

type ConfigImpl[T any] struct {
	value T
}

func (this *ConfigImpl[T]) Get() T {
	return this.value
}

func MakeConfig[T any](value T) *ConfigImpl[T] {
	return &ConfigImpl[T]{value: value}
}

type StateImpl[T any] struct {
	value    T
	onChange func(T)
}

func (this *StateImpl[T]) Get() T {
	return this.value
}

func (this *StateImpl[T]) Set(value T) {
	this.value = value
	this.onChange(value)
}

func (this *StateImpl[T]) SetOnChange(value func(T)) {
	this.onChange = value
}

func MakeState[T any]() *StateImpl[T] {
	return &StateImpl[T]{}
}

type ActionImpl[T any] struct {
	callback func(T)
}

func (this *ActionImpl[T]) RegisterCallback(callback func(T)) {
	this.callback = callback
}

func (this *ActionImpl[T]) HasCallback() bool {
	return this.callback != nil
}

func (this *ActionImpl[T]) Call(value T) {
	this.callback(value)
}

func MakeAction[T any]() *ActionImpl[T] {
	return &ActionImpl[T]{}
}
