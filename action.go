package libonebot

// ActionXxx 表示 OneBot 标准定义的核心动作名称.
const (
	// LibOneBot 自动处理的特殊动作
	ActionGetLatestEvents     = "get_latest_events"     // 获取最新事件列表
	ActionGetSupportedActions = "get_supported_actions" // 获取支持的动作列表

	// OneBot 元信息相关动作
	ActionGetStatus  = "get_status"  // 获取 OneBot 运行状态
	ActionGetVersion = "get_version" // 获取 OneBot 版本

	// 消息相关动作
	ActionSendMessage   = "send_message"   // 发送消息
	ActionDeleteMessage = "delete_message" // 删除消息

	// 用户相关动作
	ActionGetSelfInfo   = "get_self_info"   // 获取机器人自身信息
	ActionGetUserInfo   = "get_user_info"   // 获取用户信息
	ActionGetFriendList = "get_friend_list" // 获取好友列表

	// 群相关动作
	ActionGetGroupInfo       = "get_group_info"        // 获取群信息
	ActionGetGroupList       = "get_group_list"        // 获取群列表
	ActionGetGroupMemberInfo = "get_group_member_info" // 获取群成员信息
	ActionGetGroupMemberList = "get_group_member_list" // 获取群成员列表
)

// Action 表示一个动作名称.
type Action struct {
	Prefix     string // 动作名称前缀
	Name       string // 动作名称
	IsExtended bool   // 是否为扩展动作
}

// String 返回动作名称的字符串表示, 即动作请求中的 action 字段值.
func (a Action) String() string {
	if a.IsExtended {
		return a.Prefix + "_" + a.Name
	}
	return a.Name
}
