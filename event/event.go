package event

type Type string

const (
	TypeMessage Type = "message"
	TypeNotice  Type = "notice"
	TypeRequest Type = "request"
	TypeMeta    Type = "meta"
)

type Event struct {
	Platform   string `json:"platform"`
	SelfID     string `json:"self_id"`
	Type       Type   `json:"type"`
	DetailType string `json:"detail_type"`
}

type anyEvent interface {
	anyEventDummy()
}

func (e *Event) anyEventDummy() {}

type MessageEvent struct {
	Event
	UserID  string `json:"user_id"`
	GroupID string `json:"group_id,omitempty"`
	Message string `json:"message"`
}

type NoticeEvent struct {
	Event
}

type RequestEvent struct {
	Event
}

type MetaEvent struct {
	Event
}
