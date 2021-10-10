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
	UUID       string `json:"uuid"`        // 事件唯一标识符
	Platform   string `json:"platform"`    // OneBot 实现平台名称, 无需在构造时传入
	SelfID     string `json:"self_id"`     // 机器人自身 ID, 无需在构造时传入
	Time       int64  `json:"time"`        // 事件发生时间 (Unix 时间戳), 单位: 秒
	Type       string `json:"type"`        // 事件类型
	DetailType string `json:"detail_type"` // 事件详细类型
	SubType    string `json:"sub_type"`    // 事件子类型 (详细类型的下一级类型), 可为空
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

// 四种事件基本类型

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

// 标准消息事件

// PrivateMessageEvent 表示一个私聊消息事件.
type PrivateMessageEvent struct {
	MessageEvent
	UserID string `json:"user_id"` // 用户 ID
}

// MakePrivateMessageEvent 构造一个私聊消息事件.
func MakePrivateMessageEvent(time time.Time, messageID string, message Message, alt_message string, userID string) PrivateMessageEvent {
	return PrivateMessageEvent{
		MessageEvent: MakeMessageEvent(time, "private", messageID, message, alt_message),
		UserID:       userID,
	}
}

// GroupMessageEvent 表示一个群聊消息事件.
type GroupMessageEvent struct {
	MessageEvent
	GroupID string `json:"group_id"` // 群 ID
	UserID  string `json:"user_id"`  // 用户 ID
}

// MakeGroupMessageEvent 构造一个群聊消息事件.
func MakeGroupMessageEvent(time time.Time, messageID string, message Message, alt_message string, groupID string, userID string) GroupMessageEvent {
	return GroupMessageEvent{
		MessageEvent: MakeMessageEvent(time, "group", messageID, message, alt_message),
		GroupID:      groupID,
		UserID:       userID,
	}
}

// 标准通知事件

// GroupMemberIncreaseNoticeEvent 表示一个群成员增加通知事件.
type GroupMemberIncreaseNoticeEvent struct {
	NoticeEvent
	GroupID    string `json:"group_id"`    // 群 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

const (
	GroupMemberIncreaseNoticeEventSubTypeJoin   = "join"   // 成员主动加群
	GroupMemberIncreaseNoticeEventSubTypeInvite = "invite" // 成员被邀请入群
)

// MakeGroupMemberIncreaseNoticeEvent 构造一个群成员增加通知事件.
func MakeGroupMemberIncreaseNoticeEvent(time time.Time, groupID string, userID string, operatorID string) GroupMemberIncreaseNoticeEvent {
	return GroupMemberIncreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "group_member_increase"),
		GroupID:     groupID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// GroupMemberDecreaseNoticeEvent 表示一个群成员减少通知事件.
type GroupMemberDecreaseNoticeEvent struct {
	NoticeEvent
	GroupID    string `json:"group_id"`    // 群 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

const (
	GroupMemberDecreaseNoticeEventSubTypeLeave = "leave" // 成员主动退群
	GroupMemberDecreaseNoticeEventSubTypeKick  = "kick"  // 成员被踢出群
)

// MakeGroupMemberDecreaseNoticeEvent 构造一个群成员减少通知事件.
func MakeGroupMemberDecreaseNoticeEvent(time time.Time, groupID string, userID string, operatorID string) GroupMemberDecreaseNoticeEvent {
	return GroupMemberDecreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "group_member_decrease"),
		GroupID:     groupID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// GroupMemberBanNoticeEvent 表示一个群成员禁言通知事件.
type GroupMemberBanNoticeEvent struct {
	NoticeEvent
	GroupID    string `json:"group_id"`    // 群 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

// MakeGroupMemberBanNoticeEvent 构造一个群成员禁言通知事件.
func MakeGroupMemberBanNoticeEvent(time time.Time, groupID string, userID string, operatorID string) GroupMemberBanNoticeEvent {
	return GroupMemberBanNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "group_member_ban"),
		GroupID:     groupID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// GroupMemberUnbanNoticeEvent 表示一个群成员解除禁言通知事件.
type GroupMemberUnbanNoticeEvent struct {
	NoticeEvent
	GroupID    string `json:"group_id"`    // 群 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

// MakeGroupMemberUnbanNoticeEvent 构造一个群成员解除禁言通知事件.
func MakeGroupMemberUnbanNoticeEvent(time time.Time, groupID string, userID string, operatorID string) GroupMemberUnbanNoticeEvent {
	return GroupMemberUnbanNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "group_member_unban"),
		GroupID:     groupID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// GroupAdminSetNoticeEvent 表示一个群管理员设置通知事件.
type GroupAdminSetNoticeEvent struct {
	NoticeEvent
	GroupID    string `json:"group_id"`    // 群 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

// MakeGroupAdminSetNoticeEvent 构造一个群管理员设置通知事件.
func MakeGroupAdminSetNoticeEvent(time time.Time, groupID string, userID string, operatorID string) GroupAdminSetNoticeEvent {
	return GroupAdminSetNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "group_admin_set"),
		GroupID:     groupID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// GroupAdminUnsetNoticeEvent 表示一个群管理员取消通知事件.
type GroupAdminUnsetNoticeEvent struct {
	NoticeEvent
	GroupID    string `json:"group_id"`    // 群 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

// MakeGroupAdminUnsetNoticeEvent 构造一个群管理员取消通知事件.
func MakeGroupAdminUnsetNoticeEvent(time time.Time, groupID string, userID string, operatorID string) GroupAdminUnsetNoticeEvent {
	return GroupAdminUnsetNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "group_admin_unset"),
		GroupID:     groupID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// GroupMessageDeleteNoticeEvent 表示一个群消息删除通知事件.
type GroupMessageDeleteNoticeEvent struct {
	NoticeEvent
	GroupID    string `json:"group_id"`    // 群 ID
	MessageID  string `json:"message_id"`  // 消息 ID
	UserID     string `json:"user_id"`     // 消息发送者 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

const (
	GroupMessageDeleteNoticeEventSubTypeRecall = "recall" // 发送者主动撤回消息
	GroupMessageDeleteNoticeEventSubTypeDelete = "delete" // 管理员删除消息
)

// MakeGroupMessageDeleteNoticeEvent 构造一个群消息删除通知事件.
func MakeGroupMessageDeleteNoticeEvent(time time.Time, groupID string, messageID string, userID string, operatorID string) GroupMessageDeleteNoticeEvent {
	return GroupMessageDeleteNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "group_message_delete"),
		GroupID:     groupID,
		MessageID:   messageID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// FriendIncreaseNoticeEvent 表示一个好友增加通知事件.
type FriendIncreaseNoticeEvent struct {
	NoticeEvent
	UserID string `json:"user_id"` // 用户 ID
}

// MakeFriendIncreaseNoticeEvent 构造一个好友增加通知事件.
func MakeFriendIncreaseNoticeEvent(time time.Time, userID string) FriendIncreaseNoticeEvent {
	return FriendIncreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "friend_increase"),
		UserID:      userID,
	}
}

// FriendDecreaseNoticeEvent 表示一个好友减少通知事件.
type FriendDecreaseNoticeEvent struct {
	NoticeEvent
	UserID string `json:"user_id"` // 用户 ID
}

// MakeFriendDecreaseNoticeEvent 构造一个好友减少通知事件.
func MakeFriendDecreaseNoticeEvent(time time.Time, userID string) FriendDecreaseNoticeEvent {
	return FriendDecreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "friend_decrease"),
		UserID:      userID,
	}
}

// PrivateMessageDeleteNoticeEvent 表示一个私聊消息删除通知事件.
type PrivateMessageDeleteNoticeEvent struct {
	NoticeEvent
	MessageID string `json:"message_id"` // 消息 ID
	UserID    string `json:"user_id"`    // 消息发送者 ID
}

// MakePrivateMessageDeleteNoticeEvent 构造一个私聊消息删除通知事件.
func MakePrivateMessageDeleteNoticeEvent(time time.Time, messageID string, userID string) PrivateMessageDeleteNoticeEvent {
	return PrivateMessageDeleteNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "private_message_delete"),
		MessageID:   messageID,
		UserID:      userID,
	}
}

// 标准元事件

// HeartbeatMetaEvent 表示一个心跳元事件.
type HeartbeatMetaEvent struct {
	MetaEvent
	Interval int64       `json:"interval"` // 到下次心跳的间隔，单位: 秒
	Status   interface{} `json:"status"`   // OneBot 状态, 与 get_status 动作响应数据一致
}

// MakeHeartbeatMetaEvent 构造一个心跳元事件.
func MakeHeartbeatMetaEvent(time time.Time, interval int64, status interface{}) HeartbeatMetaEvent {
	return HeartbeatMetaEvent{
		MetaEvent: MakeMetaEvent(time, "heartbeat"),
		Interval:  interval,
		Status:    status,
	}
}
