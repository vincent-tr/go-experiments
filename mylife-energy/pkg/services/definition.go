package services

type Service interface {
	Init() error
	Terminate() error

	ServiceName() string
	Dependencies() []string
}
