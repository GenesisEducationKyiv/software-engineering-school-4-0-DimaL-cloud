package configs

type Rate struct {
	APIUrls struct {
		Nbu        string `yaml:"nbu"`
		PrivatBank string `yaml:"privatbank"`
	} `yaml:"api_urls"`
	NotificationCron string `yaml:"notification_cron"`
}
