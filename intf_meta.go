// 元接口

package libonebot

import "time"

// 元事件

// HeartbeatMetaEvent 表示一个心跳元事件.
type HeartbeatMetaEvent struct {
	MetaEvent
	Interval int64       `json:"interval"` // 到下次心跳的间隔，单位: 毫秒
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

// 元动作

const (
	// LibOneBot 自动处理的特殊元动作
	ActionGetLatestEvents     = "get_latest_events"     // 获取最新事件列表 (仅 HTTP 通信方式支持)
	ActionGetSupportedActions = "get_supported_actions" // 获取支持的动作列表

	ActionGetStatus  = "get_status"  // 获取 OneBot 运行状态
	ActionGetVersion = "get_version" // 获取 OneBot 版本信息
)
