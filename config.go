package libonebot

// Config 表示一个 OneBot 配置.
type Config struct {
	Heartbeat struct {
		Enabled  bool   `mapstructure:"enabled"`  // 是否启用
		Interval uint32 `mapstructure:"interval"` // 心跳间隔, 单位: 秒, 必须大于 0
	} `mapstructure:"heartbeat"` // 心跳

	Auth struct {
		AccessToken string `mapstructure:"access_token"` // 访问令牌
	} `mapstructure:"auth"` // 鉴权

	CommMethods struct {
		HTTP        []ConfigCommHTTP        `mapstructure:"http"`         // HTTP 通信方式
		HTTPWebhook []ConfigCommHTTPWebhook `mapstructure:"http_webhook"` // HTTP Webhook 通信方式
		WS          []ConfigCommWS          `mapstructure:"ws"`           // WebSocket 通信方式
		WSReverse   []ConfigCommWSReverse   `mapstructure:"ws_reverse"`   // 反向 WebSocket 通信方式
	} `mapstructure:"comm_methods"` // 通信方式
}

// ConfigCommHTTP 配置一个 HTTP 通信方式.
type ConfigCommHTTP struct {
	Host            string `mapstructure:"host"`              // HTTP 服务器监听 IP
	Port            uint16 `mapstructure:"port"`              // HTTP 服务器监听端口
	EventEnabled    bool   `mapstructure:"event_enabled"`     // 是否启用 get_latest_events 轮询动作
	EventBufferSize uint32 `mapstructure:"event_buffer_size"` // 事件缓冲区大小, 超过该大小将会丢弃最旧的事件, 0 表示不限大小
}

// ConfigCommHTTPWebhook 配置一个 HTTP Webhook 通信方式.
type ConfigCommHTTPWebhook struct {
	URL     string `mapstructure:"url"`     // Webhook 上报地址
	Timeout uint32 `mapstructure:"timeout"` // 上报请求超时时间
	Secret  string `mapstructure:"secret"`  // 签名密钥
}

// ConfigCommWS 配置一个 WebSocket 通信方式.
type ConfigCommWS struct {
	Host string `mapstructure:"host"` // WebSocket 服务器监听 IP
	Port uint16 `mapstructure:"port"` // WebSocket 服务器监听端口
}

// ConfigCommWSReverse 配置一个反向 WebSocket 通信方式.
type ConfigCommWSReverse struct {
	URL               string `mapstructure:"url"`                // 反向 WebSocket 连接地址
	ReconnectInterval uint32 `mapstructure:"reconnect_interval"` // 反向 WebSocket 重连间隔
}
