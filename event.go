package libonebot

import (
	"encoding/json"
	"time"
)

type eventType struct{ string }

func (t eventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.string)
}

// EventTypeXXX 表示 OneBot 标准定义的事件类型.
var (
	EventTypeMessage = eventType{"message"} // 消息事件
	EventTypeNotice  = eventType{"notice"}  // 通知事件
	EventTypeRequest = eventType{"request"} // 请求事件
	EventTypeMeta    = eventType{"meta"}    // 元事件
)

// Event 包含所有类型事件的共同字段.
type Event struct {
	// lock       sync.RWMutex
	Platform   string    `json:"platform"`    // OneBot 实现平台名称
	Time       int64     `json:"time"`        // 事件发生时间
	SelfID     string    `json:"self_id"`     // 机器人自身 ID
	Type       eventType `json:"type"`        // 事件类型
	DetailType string    `json:"detail_type"` // 事件详细类型
}

// AnyEvent 是所有事件对象共同实现的接口.
type AnyEvent interface {
	Name() string
	tryFixUp(platform string) bool
}

// Name 返回事件名称.
func (e *Event) Name() string {
	// e.lock.RLock()
	// defer e.lock.RUnlock()
	return e.Type.string + "." + e.DetailType
}

func (e *Event) tryFixUp(platform string) bool {
	// e.lock.Lock()
	// defer e.lock.Unlock()
	if e.SelfID == "" || e.Type.string == "" || e.DetailType == "" {
		return false
	}
	if e.Time == 0 {
		e.Time = time.Now().Unix()
	}
	e.Platform = platform
	return true
}

// MessageEvent 表示一个消息事件.
type MessageEvent struct {
	Event
	UserID  string  `json:"user_id"`            // 用户 ID
	GroupID string  `json:"group_id,omitempty"` // 群 ID
	Message Message `json:"message"`            // 消息内容
}

// NoticeEvent 表示一个通知事件.
type NoticeEvent struct {
	Event
}

// RequestEvent 表示一个请求事件.
type RequestEvent struct {
	Event
}

// MetaEvent 表示一个元事件.
type MetaEvent struct {
	Event
}
