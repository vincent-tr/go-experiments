package internal

type Plugin struct {
	Name        string `annotation:"name=name"`
	Description string `annotation:"name=description"`
	Usage       string `annotation:"name=usage,required,oneOf=sensor;actuator;logic;ui"`
}

type State struct {
	Name        string `annotation:"name=name"`
	Description string `annotation:"name=description"`
	Type        string `annotation:"name=type"`
}

type Action struct {
	Name        string `annotation:"name=name"`
	Description string `annotation:"name=description"`
	Type        string `annotation:"name=type"`
}

type Config struct {
	Name        string `annotation:"name=name"`
	Description string `annotation:"name=description"`
}
