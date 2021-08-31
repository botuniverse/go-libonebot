package libonebot

type CoreAction struct {
	name string
}

var (
	// Special
	actionGetLatestEvents = CoreAction{"get_latest_events"}

	// Meta Info
	ActionGetStatus  = CoreAction{"get_status"}
	ActionGetVersion = CoreAction{"get_version"}

	// Message
	ActionSendMessage   = CoreAction{"send_message"}
	ActionDeleteMessage = CoreAction{"delete_message"}

	// User
	ActionGetSelfInfo   = CoreAction{"get_self_info"}
	ActionGetUserInfo   = CoreAction{"get_user_info"}
	ActionGetFriendList = CoreAction{"get_friend_list"}

	// Group
	ActionGetGroupInfo       = CoreAction{"get_group_info"}
	ActionGetGroupList       = CoreAction{"get_group_list"}
	ActionGetGroupMemberInfo = CoreAction{"get_group_member_info"}
	ActionGetGroupMemberList = CoreAction{"get_group_member_list"}
)

type Action struct {
	Prefix     string
	Name       string
	IsExtended bool
}

func (a Action) String() string {
	if a.IsExtended {
		return a.Prefix + "_" + a.Name
	}
	return a.Name
}
