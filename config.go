package onebot

type Config struct {
	Heartbeat struct {
		Enabled  bool   `mapstructure:"enabled"`
		Interval uint32 `mapstructure:"interval"`
	} `mapstructure:"heartbeat"`

	Auth struct {
		AccessToken string `mapstructure:"access_token"`
	} `mapstructure:"auth"`

	CommMethods struct {
		HTTP []struct {
			Host string `mapstructure:"host"`
			Port uint16 `mapstructure:"port"`
		} `mapstructure:"http"`

		HTTPWebhook []struct {
			URL     string `mapstructure:"url"`
			Timeout uint32 `mapstructure:"timeout"`
			Secret  string `mapstructure:"secret"`
		} `mapstructure:"http_webhook"`

		WS []struct {
			Host string `mapstructure:"host"`
			Port uint16 `mapstructure:"port"`
		} `mapstructure:"ws"`

		WSReverse []struct {
			URL               string `mapstructure:"url"`
			ReconnectInterval uint32 `mapstructure:"reconnect_interval"`
		} `mapstructure:"ws_reverse"`
	} `mapstructure:"comm_methods"`
}
