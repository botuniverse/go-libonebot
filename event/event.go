package event

import (
	"sync"
	"time"
)

type type_ struct{ string }

var (
	TypeMessage type_ = type_{"message"}
	TypeNotice  type_ = type_{"notice"}
	TypeRequest type_ = type_{"request"}
	TypeMeta    type_ = type_{"meta"}
)

type Event struct {
	lock       sync.Mutex
	Platform   string `json:"platform"`
	Time       int64  `json:"time"`
	SelfID     string `json:"self_id"`
	Type       type_  `json:"type"`
	DetailType string `json:"detail_type"`
}

type AnyEvent interface {
	TryFixUp() bool
}

func (e *Event) TryFixUp() bool {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.Platform == "" || e.SelfID == "" || e.Type.string == "" || e.DetailType == "" {
		return false
	}
	if e.Time == 0 {
		e.Time = time.Now().Unix()
	}
	return true
}

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
