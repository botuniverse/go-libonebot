// OneBot Connect - 数据协议 - 基本数据类型 - 机器人自身标识
// https://12.onebot.dev/connect/data-protocol/basic-types/#_10

package libonebot

// Self 用于唯一标识一个机器人账号.
type Self struct {
	Platform string `json:"platform"`                       // 机器人平台名称
	UserID   string `json:"user_id" mapstructure:"user_id"` // 机器人自身 ID
}
