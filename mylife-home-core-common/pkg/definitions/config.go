package definitions

type Config[T any] interface {
	Get() T
}
