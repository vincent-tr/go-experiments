package services

type Service interface {
	Init(arg interface{}) error
	Terminate() error

	ServiceName() string
	Dependencies() []string
}
