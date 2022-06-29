// 单级群组接口

package libonebot

import "time"

// 群消息事件

// GroupMessageEvent 表示一个群消息事件.
type GroupMessageEvent struct {
	MessageEvent
	GroupID string `json:"group_id"` // 群 ID
	UserID  string `json:"user_id"`  // 用户 ID
}

// MakeGroupMessageEvent 构造一个群消息事件.
func MakeGroupMessageEvent(time time.Time, messageID string, message Message, alt_message string, groupID string, userID string) GroupMessageEvent {
	return GroupMessageEvent{
		MessageEvent: MakeMessageEvent(time, "group", messageID, message, alt_message),
		GroupID:      groupID,
		UserID:       userID,
	}
}

// 群通知事件

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

// 群动作

const (
	ActionGetGroupInfo       = "get_group_info"        // 获取群信息
	ActionGetGroupList       = "get_group_list"        // 获取群列表
	ActionGetGroupMemberInfo = "get_group_member_info" // 获取群成员信息
	ActionGetGroupMemberList = "get_group_member_list" // 获取群成员列表
	ActionSetGroupName       = "set_group_name"        // 设置群名称
	ActionLeaveGroup         = "leave_group"           // 退出群
)
