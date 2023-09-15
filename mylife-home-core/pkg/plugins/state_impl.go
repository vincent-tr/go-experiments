package plugins

import (
	"fmt"
	"mylife-home-core-common/metadata"
	"sync"
)

type untypedState interface {
	UntypedGet() any
	SetOnChange(value func(any))
}

type stateImpl[T comparable] struct {
	mutex    sync.Mutex
	value    T
	onChange func(any)
}

// State impl

func (state *stateImpl[T]) Get() T {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	return state.value
}

func (state *stateImpl[T]) Set(value T) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	if state.value != value {
		state.value = value
		state.onChange(value)
	}
}

// untypedState impl

func (state *stateImpl[T]) UntypedGet() any {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	return state.value
}

func (state *stateImpl[T]) SetOnChange(value func(any)) {
	state.onChange = value
}

// ---

func makeStateImpl(typ metadata.Type) untypedState {
	switch typ.(type) {
	case *metadata.RangeType:
		return &stateImpl[int64]{}
	case *metadata.TextType:
		return &stateImpl[string]{}
	case *metadata.FloatType:
		return &stateImpl[float64]{}
	case *metadata.BoolType:
		return &stateImpl[bool]{}
	case *metadata.EnumType:
		return &stateImpl[string]{}
	case *metadata.ComplexType:
		return &stateImpl[any]{}
	default:
		panic(fmt.Sprintf("Unexpected type '%s'", typ.String()))
	}
}
