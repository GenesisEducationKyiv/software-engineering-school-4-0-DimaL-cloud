package configs

type Rate struct {
	APIUrls struct {
		Nbu        string `yaml:"nbu"`
		PrivatBank string `yaml:"privatbank"`
		Fawazahmed string `yaml:"fawazahmed"`
	} `yaml:"api_urls"`
	NotificationCron string `yaml:"notification_cron"`
}
