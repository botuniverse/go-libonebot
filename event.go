package libonebot

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// EventTypeXxx 表示 OneBot 标准定义的事件类型.
const (
	EventTypeMessage = "message" // 消息事件
	EventTypeNotice  = "notice"  // 通知事件
	EventTypeRequest = "request" // 请求事件
	EventTypeMeta    = "meta"    // 元事件
)

// Event 包含所有类型事件的共同字段.
type Event struct {
	// lock       sync.RWMutex
	UUID       string `json:"uuid"`        // 事件唯一标识符
	Platform   string `json:"platform"`    // OneBot 实现平台名称, 无需在构造时传入
	SelfID     string `json:"self_id"`     // 机器人自身 ID, 无需在构造时传入
	Time       int64  `json:"time"`        // 事件发生时间, 可选, 若不传入则使用当前时间
	Type       string `json:"type"`        // 事件类型
	DetailType string `json:"detail_type"` // 事件详细类型
}

func makeEvent(time time.Time, type_ string, detailType string) Event {
	return Event{
		UUID:       uuid.New().String(),
		Time:       time.Unix(),
		Type:       type_,
		DetailType: detailType,
	}
}

// AnyEvent 是所有事件对象共同实现的接口.
type AnyEvent interface {
	Name() string
	tryFixUp(platform string, selfID string) error
	encode() ([]byte, error)
}

// Name 返回事件名称.
func (e *Event) Name() string {
	// e.lock.RLock()
	// defer e.lock.RUnlock()
	return e.Type + "." + e.DetailType
}

func (e *Event) tryFixUp(platform string, selfID string) error {
	// e.lock.Lock()
	// defer e.lock.Unlock()
	if e.Time == 0 {
		return errors.New("`time` 字段值无效")
	}
	if e.Type != EventTypeMessage && e.Type != EventTypeNotice && e.Type != EventTypeRequest && e.Type != EventTypeMeta {
		return errors.New("`type` 字段值无效")
	}
	if e.DetailType == "" {
		return errors.New("`detail_type` 字段值无效")
	}
	e.Platform = platform // override platform field directly
	if e.SelfID == "" {
		e.SelfID = selfID
	}
	return nil
}

func (e *Event) encode() ([]byte, error) {
	return json.Marshal(e)
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
func MakeHeartbeatMetaEvent(time time.Time) HeartbeatMetaEvent {
	return HeartbeatMetaEvent{
		MetaEvent: MakeMetaEvent(time, "heartbeat"),
	}
}
