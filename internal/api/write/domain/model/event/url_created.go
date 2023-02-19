package event

type UrlCreated struct {
	Model
}

func (u UrlCreated) GetEvent() UrlCreated {
	return u
}

func NewUrlCreate(data []byte) Event {
	return &UrlCreated{Model{
		EventName: "url.created",
		Data:      data,
	}}
}
