package libonebot

// ActionXxx 表示 OneBot 标准定义的核心动作名称.
const (
	// LibOneBot 自动处理的特殊动作
	ActionGetLatestEvents     = "get_latest_events"     // 获取最新事件列表 (仅 HTTP 通信方式支持)
	ActionGetSupportedActions = "get_supported_actions" // 获取支持的动作列表

	// OneBot 元信息相关动作
	ActionGetStatus  = "get_status"  // 获取 OneBot 运行状态
	ActionGetVersion = "get_version" // 获取 OneBot 版本信息

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
	ActionSetGroupName       = "set_group_name"        // 设置群名称
	ActionLeaveGroup         = "leave_group"           // 退出群
	ActionKickGroupMember    = "kick_group_member"     // 踢出群成员
	ActionBanGroupMember     = "ban_group_member"      // 禁言群成员
	ActionUnbanGroupMember   = "unban_group_member"    // 解除禁言群成员
	ActionSetGroupAdmin      = "set_group_admin"       // 设置群管理员
	ActionUnsetGroupAdmin    = "unset_group_admin"     // 取消群管理员

	// 文件相关动作
	ActionUploadFile           = "upload_file"            // 上传文件
	ActionUploadFileFragmented = "upload_file_fragmented" // 分片上传文件
	ActionGetFile              = "get_file"               // 获取文件
	ActionGetFileFragmented    = "get_file_fragmented"    // 分片获取文件
)
