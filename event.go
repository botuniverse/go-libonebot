package libonebot

import (
	"encoding/json"
	"errors"
	"time"
)

type eventType struct{ string }

func (t eventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.string)
}

// EventTypeXxx 表示 OneBot 标准定义的事件类型.
var (
	EventTypeMessage = eventType{"message"} // 消息事件
	EventTypeNotice  = eventType{"notice"}  // 通知事件
	EventTypeRequest = eventType{"request"} // 请求事件
	EventTypeMeta    = eventType{"meta"}    // 元事件
)

// Event 包含所有类型事件的共同字段.
type Event struct {
	// lock       sync.RWMutex
	Platform   string    `json:"platform"`    // OneBot 实现平台名称, 无需在构造时传入
	SelfID     string    `json:"self_id"`     // 机器人自身 ID
	Time       int64     `json:"time"`        // 事件发生时间, 可选, 若不传入则使用当前时间
	Type       eventType `json:"type"`        // 事件类型
	DetailType string    `json:"detail_type"` // 事件详细类型
}

// AnyEvent 是所有事件对象共同实现的接口.
type AnyEvent interface {
	Name() string
	tryFixUp(platform string, selfID string) error
}

// Name 返回事件名称.
func (e *Event) Name() string {
	// e.lock.RLock()
	// defer e.lock.RUnlock()
	return e.Type.string + "." + e.DetailType
}

func (e *Event) tryFixUp(platform string, selfID string) error {
	// e.lock.Lock()
	// defer e.lock.Unlock()
	if e.Time == 0 {
		return errors.New("事件 `time` 字段值无效")
	}
	if e.Type.string == "" {
		return errors.New("事件 `type` 字段值无效")
	}
	if e.DetailType == "" {
		return errors.New("事件 `detail_type` 字段值无效")
	}
	e.Platform = platform // override platform field directly
	if e.SelfID == "" {
		e.SelfID = selfID
	}
	return nil
}

// MessageEvent 表示一个消息事件.
type MessageEvent struct {
	Event
	Message Message `json:"message"` // 消息内容
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

func MakeMetaEvent(time time.Time, detailType string) MetaEvent {
	if detailType == "" {
		panic("`detail_type` 不可以为空")
	}
	return MetaEvent{
		Event: Event{
			Time:       time.Unix(),
			Type:       EventTypeMeta,
			DetailType: detailType,
		},
	}
}
