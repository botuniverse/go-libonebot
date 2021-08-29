package libonebot

type Config struct {
	Heartbeat struct {
		Enabled  bool   `mapstructure:"enabled"`
		Interval uint32 `mapstructure:"interval"`
	} `mapstructure:"heartbeat"`

	Auth struct {
		AccessToken string `mapstructure:"access_token"`
	} `mapstructure:"auth"`

	CommMethods struct {
		HTTP        []ConfigCommHTTP        `mapstructure:"http"`
		HTTPWebhook []ConfigCommHTTPWebhook `mapstructure:"http_webhook"`
		WS          []ConfigCommWS          `mapstructure:"ws"`
		WSReverse   []ConfigCommWSReverse   `mapstructure:"ws_reverse"`
	} `mapstructure:"comm_methods"`
}

type ConfigCommHTTP struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

type ConfigCommHTTPWebhook struct {
	URL     string `mapstructure:"url"`
	Timeout uint32 `mapstructure:"timeout"`
	Secret  string `mapstructure:"secret"`
}

type ConfigCommWS struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

type ConfigCommWSReverse struct {
	URL               string `mapstructure:"url"`
	ReconnectInterval uint32 `mapstructure:"reconnect_interval"`
}
