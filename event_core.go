package libonebot

import "time"

// 核心消息事件

// PrivateMessageEvent 表示一个私聊消息事件.
type PrivateMessageEvent struct {
	MessageEvent
	UserID string `json:"user_id"` // 用户 ID
}

// GroupMessageEvent 表示一个群聊消息事件.
type GroupMessageEvent struct {
	MessageEvent
	UserID  string `json:"user_id"`  // 用户 ID
	GroupID string `json:"group_id"` // 群 ID
}

// 核心元事件

// HeartbeatMetaEvent 表示一个心跳元事件.
type HeartbeatMetaEvent struct {
	MetaEvent
}

// MakeHeartbeatMetaEvent 构造一个心跳元事件.
func MakeHeartbeatMetaEvent() HeartbeatMetaEvent {
	return HeartbeatMetaEvent{
		MetaEvent: MakeMetaEvent(time.Now(), "heartbeat"),
	}
}
