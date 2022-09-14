// OneBot Connect - 数据协议 - 事件
// https://12.onebot.dev/connect/data-protocol/event/

package libonebot

import (
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
	ID         string  `json:"id"`             // 事件 ID, 构造时自动生成
	Impl       string  `json:"impl"`           // OneBot 实现名称, 无需在构造时传入
	Time       float64 `json:"time"`           // 事件发生时间 (Unix 时间戳), 单位: 秒
	Type       string  `json:"type"`           // 事件类型
	DetailType string  `json:"detail_type"`    // 事件详细类型
	SubType    string  `json:"sub_type"`       // 事件子类型 (详细类型的下一级类型), 可为空
	Self       *Self   `json:"self,omitempty"` // 机器人自身标识, 仅用于非元事件, 无需在构造时传入
}

func makeEvent(time time.Time, type_ string, detailType string) Event {
	return Event{
		ID:         uuid.New().String(),
		Time:       float64(time.UnixMicro()) / 1e6,
		Type:       type_,
		DetailType: detailType,
	}
}

// AnyEvent 是所有事件对象共同实现的接口.
type AnyEvent interface {
	Name() string
	tryFixUp(impl string, self *Self) error
}

// Name 返回事件名称.
func (e *Event) Name() string {
	// e.lock.RLock()
	// defer e.lock.RUnlock()
	return e.Type + "." + e.DetailType
}

func (e *Event) tryFixUp(impl string, self *Self) error {
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
	if e.Type == EventTypeMeta {
		if e.Self != nil {
			return errors.New("元事件中不应包含 `self` 字段")
		}
	} else {
		if self != nil {
			// prefer `self` passed in as argument
			e.Self = self
		} else if e.Self == nil {
			return errors.New("非元事件中必须包含 `self` 字段")
		}
	}
	e.Impl = impl // overwrite impl field directly
	return nil
}

// 四种事件基本类型

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

// MessageEvent 表示一个消息事件.
type MessageEvent struct {
	Event
	MessageID  string  `json:"message_id"`  // 消息 ID
	Message    Message `json:"message"`     // 消息内容
	AltMessage string  `json:"alt_message"` // 消息内容的替代表示, 可为空
}

// MakeMessageEvent 构造一个消息事件.
func MakeMessageEvent(time time.Time, detailType string, messageID string, message Message, alt_message string) MessageEvent {
	return MessageEvent{
		Event:      makeEvent(time, EventTypeMessage, detailType),
		MessageID:  messageID,
		Message:    message,
		AltMessage: alt_message,
	}
}

// NoticeEvent 表示一个通知事件.
type NoticeEvent struct {
	Event
}

// MakeNoticeEvent 构造一个通知事件.
func MakeNoticeEvent(time time.Time, detailType string) NoticeEvent {
	return NoticeEvent{
		Event: makeEvent(time, EventTypeNotice, detailType),
	}
}

// RequestEvent 表示一个请求事件.
type RequestEvent struct {
	Event
}

// MakeRequestEvent 构造一个请求事件.
func MakeRequestEvent(time time.Time, detailType string) RequestEvent {
	return RequestEvent{
		Event: makeEvent(time, EventTypeRequest, detailType),
	}
}
