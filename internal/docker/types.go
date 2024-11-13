package docker

type Container struct {
	ID      string            `json:"id"`
	Names   []string          `json:"names"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	ImageID string            `json:"imageId"`
	Command string            `json:"command"`
	Created int64             `json:"created"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	Health  string            `json:"health,omitempty"`
	Host    string            `json:"host,omitempty"`
	Tty     bool              `json:"-"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type ContainerEvent struct {
	ActorID string `json:"actorId"`
	Action  string `json:"action"`
	Host    string `json:"host"`
}

type LogEvent struct {
	Id        uint32   `json:"id,omitempty"`
	Message   any      `json:"m,omitempty"`
	Timestamp int64    `json:"ts"`
	StdType   string   `json:"t,omitempty"`
	Level     LogLevel `json:"l,omitempty"`
}
