package libonebot

// CoreAction 表示一个 OneBot 核心动作名称.
//
// 用户不应该自行创建 CoreAction 变量, 而应该使用
// ActionGetStatus, ActionGetVersion 等预定义的值.
type CoreAction struct {
	name string
}

// LibOneBot 处理的特殊动作.
var (
	actionGetLatestEvents = CoreAction{"get_latest_events"} // 获取最新事件
)

// ActionXxx 表示 OneBot 标准定义的核心动作.
var (
	// OneBot 元信息相关动作
	ActionGetStatus  = CoreAction{"get_status"}  // 获取 OneBot 运行状态
	ActionGetVersion = CoreAction{"get_version"} // 获取 OneBot 版本

	// 消息相关动作
	ActionSendMessage   = CoreAction{"send_message"}   // 发送消息
	ActionDeleteMessage = CoreAction{"delete_message"} // 删除消息

	// 用户相关动作
	ActionGetSelfInfo   = CoreAction{"get_self_info"}   // 获取机器人自身信息
	ActionGetUserInfo   = CoreAction{"get_user_info"}   // 获取用户信息
	ActionGetFriendList = CoreAction{"get_friend_list"} // 获取好友列表

	// 群相关动作
	ActionGetGroupInfo       = CoreAction{"get_group_info"}        // 获取群信息
	ActionGetGroupList       = CoreAction{"get_group_list"}        // 获取群列表
	ActionGetGroupMemberInfo = CoreAction{"get_group_member_info"} // 获取群成员信息
	ActionGetGroupMemberList = CoreAction{"get_group_member_list"} // 获取群成员列表
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
