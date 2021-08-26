package event

const (
	TYPE_MESSAGE = "message"
	TYPE_NOTICE  = "notice"
	TYPE_REQUEST = "request"
	TYPE_META    = "meta"
)

type Event struct {
	SelfID     string
	Type       string
	DetailType string
}
