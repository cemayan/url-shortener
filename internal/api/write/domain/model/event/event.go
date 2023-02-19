package event

type Event interface {
	GetEvent() UrlCreated
}

type Model struct {
	EventName string `json:"eventName,omitempty"`
	Data      []byte `json:"data,omitempty"`
}
