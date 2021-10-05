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

func makeEvent(time time.Time, type_ eventType, detailType string) Event {
	return Event{
		Time:       time.Unix(),
		Type:       type_,
		DetailType: detailType,
	}
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

// 四种事件基本类型

// MessageEvent 表示一个消息事件.
type MessageEvent struct {
	Event
	Message Message `json:"message"` // 消息内容
}

// MakeMessageEvent 构造一个消息事件.
func MakeMessageEvent(time time.Time, detailType string, message Message) MessageEvent {
	return MessageEvent{
		Event:   makeEvent(time, EventTypeMessage, detailType),
		Message: message,
	}
}

// NoticeEvent 表示一个通知事件.
type NoticeEvent struct {
	Event
}

// MakeNoticeEvent 构造一个通知事件.
func MakeNoticeEvent(time time.Time, detailType string) MetaEvent {
	return MetaEvent{
		Event: makeEvent(time, EventTypeNotice, detailType),
	}
}

// RequestEvent 表示一个请求事件.
type RequestEvent struct {
	Event
}

// MakeRequestEvent 构造一个请求事件.
func MakeRequestEvent(time time.Time, detailType string) MetaEvent {
	return MetaEvent{
		Event: makeEvent(time, EventTypeRequest, detailType),
	}
}

// MetaEvent 表示一个元事件.
type MetaEvent struct {
	Event
}

// MakeMetaEvent 构造一个元事件.
func MakeMetaEvent(time time.Time, detailType string) MetaEvent {
	return MetaEvent{
		Event: makeEvent(time, EventTypeMeta, detailType),
	}
}

// 核心消息事件

// PrivateMessageEvent 表示一个私聊消息事件.
type PrivateMessageEvent struct {
	MessageEvent
	UserID string `json:"user_id"` // 用户 ID
}

func MakePrivateMessageEvent(time time.Time, message Message, userID string) PrivateMessageEvent {
	return PrivateMessageEvent{
		MessageEvent: MakeMessageEvent(time, "private", message),
		UserID:       userID,
	}
}

// GroupMessageEvent 表示一个群聊消息事件.
type GroupMessageEvent struct {
	MessageEvent
	UserID  string `json:"user_id"`  // 用户 ID
	GroupID string `json:"group_id"` // 群 ID
}

func MakeGroupMessageEvent(time time.Time, message Message, userID string, groupID string) GroupMessageEvent {
	return GroupMessageEvent{
		MessageEvent: MakeMessageEvent(time, "group", message),
		UserID:       userID,
		GroupID:      groupID,
	}
}

// 核心元事件

// HeartbeatMetaEvent 表示一个心跳元事件.
type HeartbeatMetaEvent struct {
	MetaEvent
}

// MakeHeartbeatMetaEvent 构造一个心跳元事件.
func MakeHeartbeatMetaEvent() HeartbeatMetaEvent { // TODO: time parameter
	return HeartbeatMetaEvent{
		MetaEvent: MakeMetaEvent(time.Now(), "heartbeat"),
	}
}
