package action

type action struct{ string }

var (
	// Meta Info
	ActionGetStatus  = action{"get_status"}
	ActionGetVersion = action{"get_version"}

	// Message
	ActionSendMessage   = action{"send_message"}
	ActionDeleteMessage = action{"delete_message"}

	// User
	ActionGetSelfInfo   = action{"get_self_info"}
	ActionGetUserInfo   = action{"get_user_info"}
	ActionGetFriendList = action{"get_friend_list"}

	// Group
	ActionGetGroupInfo       = action{"get_group_info"}
	ActionGetGroupList       = action{"get_group_list"}
	ActionGetGroupMemberInfo = action{"get_group_member_info"}
	ActionGetGroupMemberList = action{"get_group_member_list"}
)
