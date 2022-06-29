// 两级群组接口

package libonebot

import "time"

// 群组消息事件

// ChannelMessageEvent 表示一个频道消息事件.
type ChannelMessageEvent struct {
	MessageEvent
	GuildID   string `json:"guild_id"`   // 群组 ID
	ChannelID string `json:"channel_id"` // 频道 ID
	UserID    string `json:"user_id"`    // 用户 ID
}

// MakeChannelMessageEvent 构造一个频道消息事件.
func MakeChannelMessageEvent(time time.Time, messageID string, message Message, alt_message string, guildId string, channelID string, userID string) ChannelMessageEvent {
	return ChannelMessageEvent{
		MessageEvent: MakeMessageEvent(time, "channel", messageID, message, alt_message),
		GuildID:      guildId,
		ChannelID:    channelID,
		UserID:       userID,
	}
}

// 群组通知事件

// GuildMemberIncreaseNoticeEvent 表示一个群组成员增加通知事件.
type GuildMemberIncreaseNoticeEvent struct {
	NoticeEvent
	GuildID    string `json:"guild_id"`    // 群组 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

const (
	GuildMemberIncreaseNoticeEventSubTypeJoin   = "join"   // 成员主动加入
	GuildMemberIncreaseNoticeEventSubTypeInvite = "invite" // 成员被邀请加入
)

// MakeGuildMemberIncreaseNoticeEvent 构造一个群组成员增加通知事件.
func MakeGuildMemberIncreaseNoticeEvent(time time.Time, guildID string, userID string, operatorID string) GuildMemberIncreaseNoticeEvent {
	return GuildMemberIncreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "guild_member_increase"),
		GuildID:     guildID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// GuildMemberDecreaseNoticeEvent 表示一个群组成员减少通知事件.
type GuildMemberDecreaseNoticeEvent struct {
	NoticeEvent
	GuildID    string `json:"guild_id"`    // 群组 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

const (
	GuildMemberDecreaseNoticeEventSubTypeLeave = "leave" // 成员主动退出
	GuildMemberDecreaseNoticeEventSubTypeKick  = "kick"  // 成员被踢出
)

// MakeGuildMemberDecreaseNoticeEvent 构造一个群组成员减少通知事件.
func MakeGuildMemberDecreaseNoticeEvent(time time.Time, guildID string, userID string, operatorID string) GuildMemberDecreaseNoticeEvent {
	return GuildMemberDecreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "guild_member_decrease"),
		GuildID:     guildID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// ChannelMemberIncreaseNoticeEvent 表示一个频道成员增加通知事件.
type ChannelMemberIncreaseNoticeEvent struct {
	NoticeEvent
	GuildID    string `json:"guild_id"`    // 群组 ID
	ChannelID  string `json:"channel_id"`  // 频道 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

const (
	ChannelMemberIncreaseNoticeEventSubTypeJoin   = "join"   // 成员主动加入
	ChannelMemberIncreaseNoticeEventSubTypeInvite = "invite" // 成员被邀请加入
)

// MakeChannelMemberIncreaseNoticeEvent 构造一个频道成员增加通知事件.
func MakeChannelMemberIncreaseNoticeEvent(time time.Time, guildID string, channelID string, userID string, operatorID string) ChannelMemberIncreaseNoticeEvent {
	return ChannelMemberIncreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "channel_member_increase"),
		GuildID:     guildID,
		ChannelID:   channelID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// ChannelMemberDecreaseNoticeEvent 表示一个频道成员减少通知事件.
type ChannelMemberDecreaseNoticeEvent struct {
	NoticeEvent
	GuildID    string `json:"guild_id"`    // 群组 ID
	ChannelID  string `json:"channel_id"`  // 频道 ID
	UserID     string `json:"user_id"`     // 用户 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

const (
	ChannelMemberDecreaseNoticeEventSubTypeLeave = "leave" // 成员主动退出
	ChannelMemberDecreaseNoticeEventSubTypeKick  = "kick"  // 成员被踢出
)

// MakeChannelMemberDecreaseNoticeEvent 构造一个频道成员减少通知事件.
func MakeChannelMemberDecreaseNoticeEvent(time time.Time, guildID string, channelID string, userID string, operatorID string) ChannelMemberDecreaseNoticeEvent {
	return ChannelMemberDecreaseNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "channel_member_decrease"),
		GuildID:     guildID,
		ChannelID:   channelID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// ChannelMessageDeleteNoticeEvent 表示一个频道消息删除通知事件.
type ChannelMessageDeleteNoticeEvent struct {
	NoticeEvent
	GuildID    string `json:"guild_id"`    // 群组 ID
	ChannelID  string `json:"channel_id"`  // 频道 ID
	MessageID  string `json:"message_id"`  // 消息 ID
	UserID     string `json:"user_id"`     // 消息发送者 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

const (
	ChannelMessageDeleteNoticeEventSubTypeRecall = "recall" // 发送者主动删除
	ChannelMessageDeleteNoticeEventSubTypeDelete = "delete" // 管理员删除
)

// MakeChannelMessageDeleteNoticeEvent 构造一个频道消息删除通知事件.
func MakeChannelMessageDeleteNoticeEvent(time time.Time, guildID string, channelID string, messageID string, userID string, operatorID string) ChannelMessageDeleteNoticeEvent {
	return ChannelMessageDeleteNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "channel_message_delete"),
		GuildID:     guildID,
		ChannelID:   channelID,
		MessageID:   messageID,
		UserID:      userID,
		OperatorID:  operatorID,
	}
}

// ChannelCreateNoticeEvent 表示一个频道新建事件.
type ChannelCreateNoticeEvent struct {
	NoticeEvent
	GuildID    string `json:"guild_id"`    // 群组 ID
	ChannelID  string `json:"channel_id"`  // 频道 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

// MakeChannelCreateNoticeEvent 构造一个频道新建事件.
func MakeChannelCreateNoticeEvent(time time.Time, guildID string, channelID string, operatorID string) ChannelCreateNoticeEvent {
	return ChannelCreateNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "channel_create"),
		GuildID:     guildID,
		ChannelID:   channelID,
		OperatorID:  operatorID,
	}
}

// ChannelDeleteNoticeEvent 表示一个频道删除事件.
type ChannelDeleteNoticeEvent struct {
	NoticeEvent
	GuildID    string `json:"guild_id"`    // 群组 ID
	ChannelID  string `json:"channel_id"`  // 频道 ID
	OperatorID string `json:"operator_id"` // 操作者 ID
}

// MakeChannelDeleteNoticeEvent 构造一个频道删除事件.
func MakeChannelDeleteNoticeEvent(time time.Time, guildID string, channelID string, operatorID string) ChannelDeleteNoticeEvent {
	return ChannelDeleteNoticeEvent{
		NoticeEvent: MakeNoticeEvent(time, "channel_delete"),
		GuildID:     guildID,
		ChannelID:   channelID,
		OperatorID:  operatorID,
	}
}

// 群组动作

const (
	ActionGetGuildInfo         = "get_guild_info"          // 获取群组信息
	ActionGetGuildList         = "get_guild_list"          // 获取群组列表
	ActionSetGuildName         = "set_guild_name"          // 设置群组名称
	ActionGetGuildMemberInfo   = "get_guild_member_info"   // 获取群组成员信息
	ActionGetGuildMemberList   = "get_guild_member_list"   // 获取群组成员列表
	ActionLeaveGuild           = "leave_guild"             // 退出群组
	ActionGetChannelInfo       = "get_channel_info"        // 获取频道信息
	ActionGetChannelList       = "get_channel_list"        // 获取频道列表
	ActionSetChannelName       = "set_channel_name"        // 设置频道名称
	ActionGetChannelMemberInfo = "get_channel_member_info" // 获取频道成员信息
	ActionGetChannelMemberList = "get_channel_member_list" // 获取频道成员列表
	ActionLeaveChannel         = "leave_channel"           // 退出频道
)
