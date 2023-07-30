package definitions

type Action[T any] interface {
	RegisterCallback(callback func(T))
}
