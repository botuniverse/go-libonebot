package libonebot

type coreAction struct{ string }

var (
	// Meta Info
	ActionGetStatus  = coreAction{"get_status"}
	ActionGetVersion = coreAction{"get_version"}

	// Message
	ActionSendMessage   = coreAction{"send_message"}
	ActionDeleteMessage = coreAction{"delete_message"}

	// User
	ActionGetSelfInfo   = coreAction{"get_self_info"}
	ActionGetUserInfo   = coreAction{"get_user_info"}
	ActionGetFriendList = coreAction{"get_friend_list"}

	// Group
	ActionGetGroupInfo       = coreAction{"get_group_info"}
	ActionGetGroupList       = coreAction{"get_group_list"}
	ActionGetGroupMemberInfo = coreAction{"get_group_member_info"}
	ActionGetGroupMemberList = coreAction{"get_group_member_list"}
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
